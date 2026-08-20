package scanners

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/padok-team/git-secret-scanner/internal/report/secret"
	"github.com/padok-team/git-secret-scanner/internal/scan/scanners/gitleaks"
	"github.com/padok-team/git-secret-scanner/internal/scan/scanners/trufflehog"
)

type RepositoryScanner struct {
	repository string
	directory  string
}

func NewRepositoryScanner(repository string, directory string) *RepositoryScanner {
	return &RepositoryScanner{
		repository: repository,
		directory:  directory,
	}
}

func (r *RepositoryScanner) ignoredLines(secrets secret.SecretSet) (secret.SecretSet, error) {
	ignoredSet := secret.NewSet()

	for _, s := range secrets {
		ignored, err := isLineIgnored(r.directory, s.Commit, s.Path, s.Line)
		if err != nil {
			return nil, err
		}
		if ignored {
			ignoredSet.Add(s)
		}
	}

	return ignoredSet, nil
}

func (r *RepositoryScanner) ignoredFingerprints() ([]secret.Fingerprint, error) {
	raw, err := loadIgnoredFingerprints(r.directory)
	if err != nil {
		return nil, err
	}

	fps := []secret.Fingerprint{}

	for _, line := range raw {
		split := strings.SplitN(line, ":", 3)
		commit := split[0]
		path := split[1]

		line, err := strconv.Atoi(split[2])
		if err != nil {
			return nil, err
		}

		fps = append(fps, secret.NewFingerprint(r.repository, commit, path, line))
	}

	return fps, nil
}

func (r *RepositoryScanner) Scan(ctx context.Context, full bool) (secret.SecretSet, error) {
	thSecrets, err := trufflehog.Scan(ctx, r.repository, r.directory, full)
	if err != nil {
		return nil, fmt.Errorf("trufflehog scan failed for repository %q: %w", r.repository, err)
	}
	glSecrets, err := gitleaks.Scan(ctx, r.repository, r.directory, full)
	if err != nil {
		return nil, fmt.Errorf("gitleaks scan failed for repository %q: %w", r.repository, err)
	}

	foundSecrets := thSecrets.Union(glSecrets)

	ignoredLines, err := r.ignoredLines(foundSecrets)
	if err != nil {
		return nil, fmt.Errorf("failed to read ignored lines for repository %q: %w", r.repository, err)
	}

	ignoredFingerprints, err := r.ignoredFingerprints()
	if err != nil {
		return nil, fmt.Errorf("failed to read ignored fingerprints for repository %q: %w", r.repository, err)
	}

	return foundSecrets.Diff(ignoredLines).DropFingerprints(ignoredFingerprints), nil
}
