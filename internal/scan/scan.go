package scan

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/padok-team/git-secret-scanner/internal/progress"
	"github.com/padok-team/git-secret-scanner/internal/report"
	"github.com/padok-team/git-secret-scanner/internal/scan/scanners/gitleaks"
	"github.com/padok-team/git-secret-scanner/internal/scan/scanners/trufflehog"
	"github.com/padok-team/git-secret-scanner/internal/scm"
	"github.com/padok-team/git-secret-scanner/internal/scm/git"
	"github.com/padok-team/git-secret-scanner/internal/utils"
	"github.com/rs/zerolog/log"
)

type ScanType int

const (
	ScanTypeFull ScanType = iota
	ScanTypeFilesOnly
)

func (t ScanType) String() string {
	switch t {
	case ScanTypeFull:
		return "full"
	case ScanTypeFilesOnly:
		return "files-only"
	default:
		// should never be reached
		return ""
	}
}

type ReportFormat int

const (
	ReportFormatJSON ReportFormat = iota
	ReportFormatCSV
)

func (f ReportFormat) String() string {
	switch f {
	case ReportFormatJSON:
		return "json"
	case ReportFormatCSV:
		return "csv"
	default:
		// should never be reached
		return ""
	}
}

// ReportFormat must implement cobra pflag.Value interface
func (f *ReportFormat) Set(s string) error {
	switch s {
	case ReportFormatJSON.String():
		*f = ReportFormatJSON
		return nil
	case ReportFormatCSV.String():
		*f = ReportFormatCSV
		return nil
	default:
		return fmt.Errorf(
			"format must be one of \"%s\" or \"%s\"",
			ReportFormatJSON.String(),
			ReportFormatCSV.String(),
		)
	}
}

// ReportFormat must implement cobra pflag.Value interface
func (v *ReportFormat) Type() string {
	return "{json,csv}"
}

type ScanArgs struct {
	ScanType               ScanType
	ReportPath             string
	ReportFormat           ReportFormat
	FingerprintsIgnorePath string
	BaselinePath           string
	MaxConcurrency         int
	ShowProgress           bool
}

func repoScanTask(ctx context.Context, repository string, s scm.Scm, full bool, w report.ReportWriter) error {
	destination := path.Join(utils.TempDirPath(), repository)

	url := s.GitRepoUrl(repository)

	err := git.Clone(ctx, url, destination, !full, true)
	if err != nil {
		// if remote is empty, scan next repository
		if errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return nil
		}
		return fmt.Errorf("failed to clone repository %q: %w", repository, err)
	}

	thSecrets, err := trufflehog.Scan(ctx, repository, destination, full)
	if err != nil {
		return fmt.Errorf("trufflehog scan failed for repository %q: %w", repository, err)
	}
	glSecrets, err := gitleaks.Scan(ctx, repository, destination, full)
	if err != nil {
		return fmt.Errorf("gitleaks scan failed for repository %q: %w", repository, err)
	}

	secrets := thSecrets.
		Union(glSecrets).
		Diff(GetRepoBaseline(repository)).
		DropFingerprints(GetRepoIgnoredFingerprints(repository))

	if err := w.WriteAll(secrets.ToSlice()); err != nil {
		return fmt.Errorf("failed to add secrets in report for repository %q: %w", repository, err)
	}

	if secrets.Length() > 0 {
		log.Warn().
			Int("num", secrets.Length()).
			Str("repository", repository).
			Msgf("scanned %s, got findings", repository)
	} else {
		log.Debug().
			Str("repository", repository).
			Msgf("scanned %s", repository)
	}

	return nil
}

func Scan(ctx context.Context, s scm.Scm, args ScanArgs) error {
	var err error

	log.Info().Msg("scan initiated")

	var writer report.ReportWriter

	switch args.ReportFormat {
	case ReportFormatJSON:
		writer, err = report.NewJSONReportWriter(args.ReportPath)
	case ReportFormatCSV:
		writer, err = report.NewCSVReportWriter(args.ReportPath)
	}

	if err != nil {
		return fmt.Errorf("failed to initialize report writer: %w", err)
	}

	defer writer.Close()

	if err := LoadBaseline(args.BaselinePath); err != nil {
		return fmt.Errorf("failed to load baseline: %w", err)
	}
	if err := LoadIgnoredFingerprints(args.FingerprintsIgnorePath); err != nil {
		return fmt.Errorf("failed to load ignored fingerprints: %w", err)
	}

	repos, err := s.ListRepos(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve repositories list: %w", err)
	}

	tasks := make([]*progress.Task, 0, len(repos))
	for _, repo := range repos {
		repo := repo // closure
		tasks = append(tasks, progress.NewTask(func(ctx context.Context) error {
			return repoScanTask(ctx, repo, s, args.ScanType == ScanTypeFull, writer)
		}))
	}

	if args.ShowProgress {
		err = progress.RunTasksWithProgressBar(ctx, "Scanning repositories...", tasks, args.MaxConcurrency)
	} else {
		err = progress.RunTasks(ctx, tasks, args.MaxConcurrency)
	}
	if err != nil {
		return fmt.Errorf("repositories scan failed: %w", err)
	}

	log.Info().Msg("scan completed")

	return nil
}
