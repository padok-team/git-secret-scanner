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

type ScanArgs struct {
	ScanType               ScanType
	ReportPath             string
	FingerprintsIgnorePath string
	BaselinePath           string
	MaxConcurrency         int
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

	thSecrets, err := trufflehog.TrufflehogScan(ctx, repository, destination, full)
	if err != nil {
		return fmt.Errorf("trufflehog scan failed for repository %q: %w", repository, err)
	}
	glSecrets, err := gitleaks.GitleaksScan(ctx, repository, destination, full)
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

	if utils.FileExistsAndNotEmpty(args.ReportPath) {
		return fmt.Errorf("file \"%s\" already exists", args.ReportPath)
	}

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

	writer, err := report.NewCSVReportWriter(args.ReportPath, true)
	if err != nil {
		return fmt.Errorf("failed to initialize report writer: %w", err)
	}
	defer writer.Close()

	tasks := make([]*progress.Task, 0, len(repos))
	for _, repo := range repos {
		repo := repo // closure
		tasks = append(tasks, progress.NewTask(func(ctx context.Context) error {
			return repoScanTask(ctx, repo, s, args.ScanType == ScanTypeFull, writer)
		}))
	}

	if err := progress.RunTasksWithProgressBar(ctx, "Scanning repositories...", tasks, args.MaxConcurrency); err != nil {
		return fmt.Errorf("repositories scan failed: %w", err)
	}

	log.Info().Msg("scan completed")

	return nil
}
