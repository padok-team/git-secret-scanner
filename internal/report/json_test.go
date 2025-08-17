package report

import (
	"bufio"
	"encoding/json"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/padok-team/git-secret-scanner/internal/report/secret"
)

var jsonReportPath = path.Join("../../test/report_test.json")

func TestNewJSONReportWriter(t *testing.T) {
	writer, err := NewJSONReportWriter(jsonReportPath)
	if err != nil {
		t.Fatalf(`NewJSONReportWriter("test/report.json") = _, %v, want _, nil`, err)
	}

	writer.file.WriteString("line1\n") //nolint:errcheck
	writer.file.WriteString("line2\n") //nolint:errcheck
	writer.file.WriteString("line3\n") //nolint:errcheck

	f, err := os.Open(jsonReportPath)
	if err != nil {
		t.Fatalf(`os.Open("test/report.json") = _, %v, want _, nil`, err)
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
	scanner.Scan()
	if scanner.Text() != "line3" {
		t.Fatalf(`scanner.Text() = %s, want %s`, err, "line3")
	}

	writer.Close()
	f.Close()
}

func TestJSONReportWriterWriteAll(t *testing.T) {
	writer, err := NewJSONReportWriter(jsonReportPath)
	if err != nil {
		t.Fatalf(`NewJSONReportWriter("test/report.json") = %v, want false, nil`, err)
	}

	err = writer.WriteAll([]*secret.Secret{&test1, &test2, &test3})
	if err != nil {
		t.Fatalf(`writer.WriteAll({test1, test2, test3}) = %v want nil`, err)
	}

	b, err := os.ReadFile(jsonReportPath)
	if err != nil {
		t.Fatalf(`os.ReadFile("test/report.json") = _, %v, want _, nil`, err)
	}

	want, err := json.Marshal([]*secret.Secret{&test1, &test2, &test3})
	if err != nil {
		t.Fatalf(`json.Marshal = %v, want nil`, err)
	}

	if !reflect.DeepEqual(b, want) {
		t.Fatalf(`b = %v, want %v`, b, want)
	}

	writer.Close()
}

func TestReadJSONReport(t *testing.T) {
	writer, err := NewJSONReportWriter(jsonReportPath)
	if err != nil {
		t.Fatalf(`NewJSONReportWriter("test/report.json") = _, %v, want _, nil`, err)
	}

	slice := []*secret.Secret{&test1, &test2, &test3}

	err = writer.WriteAll(slice)
	if err != nil {
		t.Fatalf(`writer.WriteAll({test1, test2, test3}) = %v want nil`, err)
	}

	s, err := ReadJSONReport(jsonReportPath)

	if !reflect.DeepEqual(slice, s) || err != nil {
		t.Fatalf(`ReadJSONReport("test/report.json") = %v, %v want %v, nil`, s, err, slice)
	}
}
