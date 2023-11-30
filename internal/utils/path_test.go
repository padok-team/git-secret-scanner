package utils

import (
	"path"
	"testing"
)

func TestTempDirPath(t *testing.T) {
	t.Setenv("TMPDIR", "")

	test := TempDirPath()
	want := path.Join("/tmp", TempDirName)

	if test != want {
		t.Fatalf(`TempDirPath() = %q, want %q`, test, want)
	}

	t.Setenv("TMPDIR", "/this/is/a/test")

	test = TempDirPath()
	want = path.Join("/this/is/a/test", TempDirName)

	if test != want {
		t.Fatalf(`TempDirPath() = %q, want %q`, test, want)
	}
}

func TestCommandExists(t *testing.T) {
	test := CommandExists("ls")
	want := true

	if test != want {
		t.Fatalf(`CommandExists("ls") = %t, want %t`, test, want)
	}

	test = CommandExists("thisisnotacommand!")
	want = false

	if test != want {
		t.Fatalf(`CommandExists("thisisnotacommand!") = %t, want %t`, test, want)
	}
}

func TestFileExistsAndNotEmpty(t *testing.T) {
	test := FileExistsAndNotEmpty("../../test/testdata/baseline.csv")
	want := true

	if test != want {
		t.Fatalf(`FileExistsAndNotEmpty("test/testdata/baseline.csv") = %t, want %t`, test, want)
	}

	test = FileExistsAndNotEmpty("../../test/testdata/i/do/not/exists")
	want = false

	if test != want {
		t.Fatalf(`FileExistsAndNotEmpty("test/testdata/i/do/not/exists") = %t, want %t`, test, want)
	}
}
