package report

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/padok-team/git-secret-scanner/internal/report/secret"
	"github.com/padok-team/git-secret-scanner/internal/utils"
)

type CSVReportWriter struct {
	SecretsWrittenCount int

	file *os.File
}

func NewCSVReportWriter(path string) (*CSVReportWriter, error) {
	path = reportPath(path, "csv")

	if utils.FileExistsAndNotEmpty(path) {
		return nil, fmt.Errorf("file \"%s\" already exists", path)
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &CSVReportWriter{SecretsWrittenCount: 0, file: f}, nil
}

func (w *CSVReportWriter) WriteAll(s []*secret.Secret) error {
	if len(s) > 0 {
		if w.SecretsWrittenCount == 0 {
			if err := gocsv.Marshal(s, w.file); err != nil {
				return err
			}
		} else {
			if err := gocsv.MarshalWithoutHeaders(s, w.file); err != nil {
				return err
			}
		}

		w.SecretsWrittenCount += len(s)
	}

	return nil
}

func (w *CSVReportWriter) Close() error {
	return w.file.Close()
}

func ReadCSVReport(path string) ([]*secret.Secret, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	secrets := make([]*secret.Secret, 0)
	if err := gocsv.UnmarshalFile(f, &secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}
