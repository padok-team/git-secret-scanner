package report

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/padok-team/git-secret-scanner/internal/report/secret"
)

type CSVReportWriter struct {
	file *os.File
}

func NewCSVReportWriter(path string, forceRecreate bool) (*CSVReportWriter, error) {
	var err error
	var f *os.File

	if forceRecreate {
		f, err = os.Create(path)
	} else {
		f, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0755)
	}
	if err != nil {
		return nil, err
	}

	return &CSVReportWriter{file: f}, nil
}

func (w *CSVReportWriter) WriteAll(s []*secret.Secret) error {
	info, err := w.file.Stat()
	if err != nil {
		return err
	}

	if info.Size() > 0 {
		return gocsv.MarshalWithoutHeaders(s, w.file)
	}

	return gocsv.Marshal(s, w.file)
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
