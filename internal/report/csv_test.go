package report

import (
	"bufio"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/gocarina/gocsv"
	"github.com/padok-team/git-secret-scanner/internal/report/secret"
)

var reportPath = path.Join("../../test/report_test.csv")

var test1 secret.Secret = secret.Secret{
	Repository:  "test1_repo",
	Path:        "test1_path",
	Kind:        secret.SecretKindGeneric,
	Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
	Line:        10,
	Valid:       secret.SecretValidityInvalid,
	Cleartext:   "test1_cleartext",
	Fingerprint: "e272f5f2466d724d16b053356da74bfa90a7b1cd3dbb579314e47da4ee43afa3",
}

var test2 secret.Secret = secret.Secret{
	Repository:  "test2_repo",
	Path:        "test2_path",
	Kind:        secret.SecretKindGeneric,
	Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42599",
	Line:        42,
	Valid:       secret.SecretValidityUnknown,
	Cleartext:   "",
	Fingerprint: "hardcoded_fingerprint",
}

var test3 secret.Secret = secret.Secret{
	Repository:  "test1_repo",
	Path:        "test1_path",
	Kind:        secret.SecretKindAMQP,
	Commit:      "e2acefa38de1bb02673cd2496a4663a3c6c42508",
	Line:        0,
	Valid:       secret.SecretValidityUnknown,
	Cleartext:   "test1_cleartext",
	Fingerprint: "e272f5f2466d724d16b053356da74bfa90a7b1cd3dbb579314e47da4ee43afa3",
}

func TestNewCSVReportWriter(t *testing.T) {
	writer, err := NewCSVReportWriter(reportPath, true)
	if err != nil {
		t.Fatalf(`NewCSVReportWriter("test/report.csv", true) = _, %v, want _, nil`, err)
	}

	writer.file.WriteString("line1\n") //nolint:errcheck
	writer.file.WriteString("line2\n") //nolint:errcheck

	f, err := os.Open(reportPath)
	if err != nil {
		t.Fatalf(`os.Open(reportPath) = _, %v, want _, nil`, err)
	}

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	if scanner.Text() != "line1" {
		t.Fatalf(`scanner.Text() = %s, want %s`, err, "line1")
	}
	scanner.Scan()
	if scanner.Text() != "line2" {
		t.Fatalf(`scanner.Text() = %s, want %s`, err, "line2")
	}

	writer.Close()
	f.Close()

	writer, err = NewCSVReportWriter(reportPath, false)
	if err != nil {
		t.Fatalf(`NewCSVReportWriter("test/report.csv", true) = _, %v, want _, nil`, err)
	}

	writer.file.WriteString("line3\n") //nolint:errcheck

	f, err = os.Open(reportPath)
	if err != nil {
		t.Fatalf(`os.Open(reportPath) = _, %v, want _, nil`, err)
	}

	scanner = bufio.NewScanner(f)

	scanner.Scan()
	if scanner.Text() != "line1" {
		t.Fatalf(`scanner.Text() = %q, want %q`, err, "line1")
	}
	scanner.Scan()
	if scanner.Text() != "line2" {
		t.Fatalf(`scanner.Text() = %q, want %q`, err, "line2")
	}
	scanner.Scan()
	if scanner.Text() != "line3" {
		t.Fatalf(`scanner.Text() = %q, want %q`, err, "line3")
	}

	writer.Close()
	f.Close()
}

func TestCSVReportWriterWriteAll(t *testing.T) {
	writer, err := NewCSVReportWriter(reportPath, true)
	if err != nil {
		t.Fatalf(`NewCSVReportWriter("test/report.csv", true) = %v, want false, nil`, err)
	}

	err = writer.WriteAll([]*secret.Secret{&test1, &test2, &test3})
	if err != nil {
		t.Fatalf(`writer.WriteAll({test1, test2, test3}) = %v want nil`, err)
	}

	f, err := os.Open(reportPath)
	if err != nil {
		t.Fatalf(`os.Open(reportPath) = _, %v, want _, nil`, err)
	}

	scanner := bufio.NewScanner(f)

	csv, err := gocsv.MarshalString([]*secret.Secret{&test1, &test2, &test3})
	if err != nil {
		t.Fatalf(`gocsv.MarshalString =%v, want nil`, err)
	}

	wants := strings.Split(csv, "\n")

	scanner.Scan()
	test := scanner.Text()
	want := wants[0]

	if scanner.Text() != want {
		t.Fatalf(`scanner.Text() = %q want %q`, test, want)
	}

	scanner.Scan()
	test = scanner.Text()
	want = wants[1]

	if scanner.Text() != want || err != nil {
		t.Fatalf(`scanner.Text() = %q, %v, want %q, nil`, test, err, want)
	}

	scanner.Scan()
	test = scanner.Text()
	want = wants[2]

	if scanner.Text() != want || err != nil {
		t.Fatalf(`scanner.Text() = %q, %v, want %q, nil`, test, err, want)
	}

	scanner.Scan()
	test = scanner.Text()
	want = wants[3]

	if scanner.Text() != want || err != nil {
		t.Fatalf(`scanner.Text() = %q, %v, want %q, nil`, test, err, want)
	}

	writer.Close()
	f.Close()
}

func TestReadCSVReport(t *testing.T) {
	writer, err := NewCSVReportWriter(reportPath, true)
	if err != nil {
		t.Fatalf(`NewCSVReportWriter("test/report.csv", true) = _, %v, want _, nil`, err)
	}

	slice := []*secret.Secret{&test1, &test2, &test3}

	err = writer.WriteAll(slice)
	if err != nil {
		t.Fatalf(`writer.WriteAll({test1, test2, test3}) = %v want nil`, err)
	}

	s, err := ReadCSVReport(reportPath)

	if !reflect.DeepEqual(slice, s) || err != nil {
		t.Fatalf(`ReadCSVReport("test/report.csv") = %v, %v want %v, nil`, s, err, slice)
	}
}
