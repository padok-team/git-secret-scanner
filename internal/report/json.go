package report

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/padok-team/git-secret-scanner/internal/report/secret"
	"github.com/padok-team/git-secret-scanner/internal/utils"
)

type JSONReportWriter struct {
	SecretsWrittenCount int

	file *os.File
}

func NewJSONReportWriter(path string) (*JSONReportWriter, error) {
	path = reportPath(path, "json")

	if utils.FileExistsAndNotEmpty(path) {
		return nil, fmt.Errorf("file \"%s\" already exists", path)
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &JSONReportWriter{SecretsWrittenCount: 0, file: f}, nil
}

func (w *JSONReportWriter) WriteAll(s []*secret.Secret) error {
	if len(s) > 0 {
		for _, secret := range s {
			var input string

			// handle separation between secrets in JSON
			if w.SecretsWrittenCount == 0 {
				// add the opening "["
				input = "[\n  "
			} else {
				// move back to overwrite the last "\n]"
				if _, err := w.file.Seek(-2, 1); err != nil {
					return err
				}
				// add the "," separation between secrets
				input = ",\n  "
			}

			b, err := json.MarshalIndent(secret, "  ", "  ")
			if err != nil {
				return err
			}

			// add the closing "]"
			input += string(b) + "\n]"

			if _, err := w.file.WriteString(input); err != nil {
				return err
			}

			w.SecretsWrittenCount++
		}
	}

	return nil
}

func (w *JSONReportWriter) Close() error {
	return w.file.Close()
}

func ReadJSONReport(path string) ([]*secret.Secret, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	secrets := make([]*secret.Secret, 0)
	if err := json.Unmarshal(b, &secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}
