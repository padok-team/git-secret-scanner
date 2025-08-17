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
	defer writer.Close()

	writer.file.WriteString("line1\n") //nolint:errcheck
	writer.file.WriteString("line2\n") //nolint:errcheck
	writer.file.WriteString("line3\n") //nolint:errcheck

	f, err := os.Open(jsonReportPath)
	if err != nil {
		t.Fatalf(`os.Open("test/report.json") = _, %v, want _, nil`, err)
	}
	defer f.Close()
	defer os.Remove(jsonReportPath)

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	want := scanner.Text()

	if want != "line1" {
		t.Fatalf(`scanner.Text() = %s, want %s`, want, "line1")
	}

	scanner.Scan()
	want = scanner.Text()

	if want != "line2" {
		t.Fatalf(`scanner.Text() = %s, want %s`, want, "line2")
	}

	scanner.Scan()
	want = scanner.Text()

	if want != "line3" {
		t.Fatalf(`scanner.Text() = %s, want %s`, want, "line3")
	}
}

func TestJSONReportWriterWriteAll(t *testing.T) {
	writer, err := NewJSONReportWriter(jsonReportPath)
	if err != nil {
		t.Fatalf(`NewJSONReportWriter("test/report.json") = %v, want false, nil`, err)
	}
	defer writer.Close()
	defer os.Remove(jsonReportPath)

	err = writer.WriteAll([]*secret.Secret{&test1, &test2, &test3})
	if err != nil {
		t.Fatalf(`writer.WriteAll({test1, test2, test3}) = %v want nil`, err)
	}

	b, err := os.ReadFile(jsonReportPath)
	if err != nil {
		t.Fatalf(`os.ReadFile("test/report.json") = _, %v, want _, nil`, err)
	}

	var s []*secret.Secret
	if err := json.Unmarshal(b, &s); err != nil {
		t.Fatalf(`json.Unmarshal = %v, want nil`, err)
	}

	want := []*secret.Secret{&test1, &test2, &test3}

	if !reflect.DeepEqual(s, want) {
		t.Fatalf(`s = %v, want %v`, s, want)
	}
}

func TestReadJSONReport(t *testing.T) {
	writer, err := NewJSONReportWriter(jsonReportPath)
	if err != nil {
		t.Fatalf(`NewJSONReportWriter("test/report.json") = _, %v, want _, nil`, err)
	}
	defer writer.Close()
	defer os.Remove(jsonReportPath)

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
