package scanners

import (
	"strings"

	"github.com/padok-team/git-secret-scanner/internal/scm/git"
)

const (
	GitleaksIgnoreTag   string = "gitleaks:allow"
	TrufflehogIgnoreTag string = "trufflehog:ignore"
)

func IsLineIgnored(repo string, hash string, file string, lineNum int) (bool, error) {
	diff, err := git.CommitDiff(repo, hash, file)
	if err != nil {
		return false, err
	}

	lines := strings.Split(diff, "\n")

	if lineNum > 0 && lineNum <= len(lines) {
		line := lines[lineNum-1]
		if strings.Contains(line, GitleaksIgnoreTag) || strings.Contains(line, TrufflehogIgnoreTag) {
			return true, nil
		}
	}

	return false, nil
}
