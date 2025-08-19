package gitleaks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/padok-team/git-secret-scanner/internal/report/secret"
	"github.com/padok-team/git-secret-scanner/internal/scan/scanners"
)

type GitleaksReportItem struct {
	Commit    string `json:"Commit"`
	File      string `json:"File"`
	StartLine int    `json:"StartLine"`
	Rule      string `json:"RuleID"`
	Secret    string `json:"Secret"`
}

func (gi *GitleaksReportItem) ToSecret(repository string) (*secret.Secret, error) {
	kind, ok := GitleaksSecretKindMapping[gi.Rule]
	if !ok {
		kind = secret.SecretKindGeneric
	}
	return secret.NewSecret(
		repository,
		gi.File,
		kind,
		gi.Commit,
		gi.StartLine,
		secret.SecretValidityUnknown,
		secret.SecretScannersGitleaks,
		gi.Secret,
		"",
	)
}

func Scan(ctx context.Context, repository string, directory string, full bool) (secret.SecretSet, error) {
	reportPath := path.Join(directory, "gitleaks.json")

	args := []string{
		"git",
		"--report-path=" + reportPath,
		"--report-format=json",
		"--no-banner",
		"--no-color",
		"--log-level=error",
		"--exit-code=0",
	}
	if !full {
		args = append(args, "--log-opts=HEAD^!")
	}

	args = append(args, directory)

	cmd := exec.CommandContext(ctx, "gitleaks", args...)

	stderrP, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	stderr, err := io.ReadAll(stderrP)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return nil, fmt.Errorf("\"gitleaks\" exited with error:\n%s", stderr)
		}
		return nil, err
	}

	f, err := os.Open(reportPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	secrets := secret.NewSet()
	decoder := json.NewDecoder(f)

	report := make([]GitleaksReportItem, 0)
	if err := decoder.Decode(&report); err != nil {
		return nil, err
	}

	for _, item := range report {
		ignored, err := scanners.IsLineIgnored(directory, item.Commit, item.File, item.StartLine)
		if err != nil {
			return nil, err
		}

		if !ignored {
			s, err := item.ToSecret(repository)
			if err != nil {
				return nil, err
			}
			secrets.Add(s)
		}
	}

	return secrets, nil
}
