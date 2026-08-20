package scanners

import (
	"errors"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/padok-team/git-secret-scanner/internal/scm/git"
)

const (
	GitleaksIgnoreTag   string = "gitleaks:allow"
	TrufflehogIgnoreTag string = "trufflehog:ignore"
)

func isLineIgnored(repo string, hash string, file string, lineNum int) (bool, error) {
	c, err := git.ReadCommitFileContent(repo, hash, file)
	if err != nil {
		return false, err
	}

	lines := strings.Split(c, "\n")

	if lineNum > 0 && lineNum <= len(lines) {
		line := lines[lineNum-1]
		if strings.Contains(line, GitleaksIgnoreTag) || strings.Contains(line, TrufflehogIgnoreTag) {
			return true, nil
		}
	}

	return false, nil
}

const GitleaksIgnoreFingerprintsFileName string = ".gitleaksignore"

func loadIgnoredFingerprints(repo string) ([]string, error) {
	c, err := git.ReadFileContent(repo, GitleaksIgnoreFingerprintsFileName)
	if err != nil {
		if errors.Is(errors.Unwrap(err), object.ErrFileNotFound) {
			// If the file doesn't exist, there are no ignored fingerprints
			return []string{}, nil
		}
		return nil, err
	}

	return strings.Split(c, "\n"), nil
}
