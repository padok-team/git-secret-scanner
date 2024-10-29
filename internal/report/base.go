package report

import "github.com/padok-team/git-secret-scanner/internal/report/secret"

type ReportWriter interface {
	WriteAll(s []*secret.Secret) error
	Close() error
}
