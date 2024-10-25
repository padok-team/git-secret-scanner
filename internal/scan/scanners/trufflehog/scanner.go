package trufflehog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/padok-team/git-secret-scanner/internal/report/secret"
	"github.com/padok-team/git-secret-scanner/internal/scan/scanners"
)

type TrufflehogReportItem struct {
	Commit       string
	File         string
	Line         int
	DetectorName string
	Verified     secret.SecretValidity
	Raw          string
}

func (ri *TrufflehogReportItem) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	gitInfo := data["SourceMetadata"].(map[string]interface{})["Data"].(map[string]interface{})["Git"].(map[string]interface{})
	detectorName := data["DetectorName"].(string)

	var file string
	v, ok := gitInfo["file"]
	if ok {
		file = v.(string)
	}

	var verified secret.SecretValidity
	if detectorName != string(secret.SecretKindPrivateKey) {
		v, ok := data["Verified"]
		if ok {
			if v.(bool) {
				verified = secret.SecretValidityValid
			} else {
				verified = secret.SecretValidityInvalid
			}
		} else {
			verified = secret.SecretValidityUnknown
		}
	}

	*ri = TrufflehogReportItem{
		Commit:       gitInfo["commit"].(string),
		File:         file,
		Line:         int(gitInfo["line"].(float64)),
		DetectorName: detectorName,
		Verified:     verified,
		Raw:          data["Raw"].(string),
	}

	return nil
}

func (ri *TrufflehogReportItem) ToSecret(repository string) (*secret.Secret, error) {
	kind, ok := TrufflehogSecretKindMapping[ri.DetectorName]
	if !ok {
		kind = secret.SecretKindGeneric
	}
	return secret.NewSecret(
		repository,
		ri.File,
		kind,
		ri.Commit,
		ri.Line,
		ri.Verified,
		secret.SecretScannersTrufflehog,
		ri.Raw,
		"",
	)
}

func TrufflehogScan(ctx context.Context, repository string, directory string, full bool) (secret.SecretSet, error) {
	args := []string{
		"git",
		"--no-update",
		"--force-skip-binaries",
		"--bare",
		"--json",
	}
	if !full {
		args = append(args, "--branch=HEAD", "--max-depth=2")
	}

	args = append(args, "file://"+directory)

	cmd := exec.CommandContext(ctx, "trufflehog", args...)

	stdoutP, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrP, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	secrets := secret.NewSet()
	decoder := json.NewDecoder(stdoutP)

	for decoder.More() {
		var item TrufflehogReportItem
		if err := decoder.Decode(&item); err != nil {
			return nil, err
		}

		// if "File" is not defined, it means Trufflehog found something in the commit message
		// this is badly handled by Trufflehog at the moment, so let's skip this for now
		if item.File == "" {
			continue
		}

		ignored, err := scanners.IsLineIgnored(directory, item.Commit, item.File, item.Line)
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

	stderr, err := io.ReadAll(stderrP)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return nil, fmt.Errorf("\"trufflehog\" exited with error:\n%s", stderr)
		}
		return nil, err
	}

	return secrets, nil
}
