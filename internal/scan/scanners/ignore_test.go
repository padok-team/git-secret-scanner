package scanners

import (
	"context"
	"path"
	"testing"

	"github.com/padok-team/git-secret-scanner/internal/scm/git"
	"github.com/padok-team/git-secret-scanner/internal/utils"
)

var testPath string = path.Join(utils.TempDirPath(), "tests")

func TestIsLineIgnored(t *testing.T) {
	err := git.Clone(context.Background(), "https://github.com/gitleaks/gitleaks", path.Join(testPath, "gitleaks"), false, false, true)
	if err != nil {
		t.Fatalf(`git.Clone("https://github.com/gitleaks/gitleaks", ...) = %v, nil`, err)
	}

	test, err := IsLineIgnored(path.Join(testPath, "gitleaks"), "abfd0f3fdcb7925ff94184fba67b5d444cc42f92", "README.md", 396)
	want := true

	if test != want || err != nil {
		t.Fatalf(`IsLineIgnored(gitleaks, "abfd0f3fdcb7925ff94184fba67b5d444cc42f92", ...) = %t, %v, want %t, nil`, test, err, want)
	}

	err = git.Clone(context.Background(), "https://github.com/trufflesecurity/trufflehog", path.Join(testPath, "trufflehog"), false, false, true)
	if err != nil {
		t.Fatalf(`git.Clone("https://github.com/trufflesecurity/trufflehog", ...) = %v, nil`, err)
	}

	test, err = IsLineIgnored(path.Join(testPath, "trufflehog"), "5d7e6fc2fa98df153a7e685f4e90508d3aea8922", "pkg/engine/engine.go", 530)
	want = true

	if test != want || err != nil {
		t.Fatalf(`IsLineIgnored(trufflehog, "5d7e6fc2fa98df153a7e685f4e90508d3aea8922", ...) = %t, %v, want %t, nil`, test, err, want)
	}
}
