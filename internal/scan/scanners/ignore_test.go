package scanners

import (
	"context"
	"path"
	"reflect"
	"testing"

	"github.com/padok-team/git-secret-scanner/internal/scm/git"
	"github.com/padok-team/git-secret-scanner/internal/utils"
)

var testPath string = path.Join(utils.TempDirPath(), "tests")

func TestIsLineIgnored(t *testing.T) {
	err := git.Clone(context.Background(), "https://github.com/gitleaks/gitleaks", path.Join(testPath, "gitleaks"), false, true)
	if err != nil {
		t.Fatalf(`git.Clone("https://github.com/gitleaks/gitleaks", ...) = %v, nil`, err)
	}

	test, err := isLineIgnored(path.Join(testPath, "gitleaks"), "abfd0f3fdcb7925ff94184fba67b5d444cc42f92", "README.md", 396)
	want := true

	if test != want || err != nil {
		t.Fatalf(`isLineIgnored(gitleaks, "abfd0f3fdcb7925ff94184fba67b5d444cc42f92", ...) = %t, %v, want %t, nil`, test, err, want)
	}

	err = git.Clone(context.Background(), "https://github.com/trufflesecurity/trufflehog", path.Join(testPath, "trufflehog"), false, true)
	if err != nil {
		t.Fatalf(`git.Clone("https://github.com/trufflesecurity/trufflehog", ...) = %v, nil`, err)
	}

	test, err = isLineIgnored(path.Join(testPath, "trufflehog"), "5d7e6fc2fa98df153a7e685f4e90508d3aea8922", "pkg/engine/engine.go", 530)
	want = true

	if test != want || err != nil {
		t.Fatalf(`isLineIgnored(trufflehog, "5d7e6fc2fa98df153a7e685f4e90508d3aea8922", ...) = %t, %v, want %t, nil`, test, err, want)
	}
}

func TestLoadIgnoredFingerprints(t *testing.T) {
	err := git.Clone(context.Background(), "https://github.com/gitleaks/gitleaks", path.Join(testPath, "gitleaks"), false, true)
	if err != nil {
		t.Fatalf(`git.Clone("https://github.com/gitleaks/gitleaks", ...) = %v, nil`, err)
	}

	test, err := loadIgnoredFingerprints(path.Join(testPath, "gitleaks"))
	want := []string{
		"418edf165dbb63d6f46993ae8f8818ffd87ea582:cmd/generate/config/rules/jwt.go:jwt:17",
		"418edf165dbb63d6f46993ae8f8818ffd87ea582:cmd/generate/config/rules/jwt.go:jwt:19",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-sensitive-url:46",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-sensitive-url:48",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-sensitive-url:50",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-sensitive-url:52",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-sensitive-url:54",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-sensitive-url:55",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-sensitive-url:56",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-sensitive-url:57",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-secret:22",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-secret:23",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-secret:24",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-secret:28",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:cmd/generate/config/rules/sidekiq.go:sidekiq-secret:29",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:detect/detect_test.go:sidekiq-sensitive-url:164",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:detect/detect_test.go:sidekiq-sensitive-url:170",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:detect/detect_test.go:sidekiq-secret:120",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:detect/detect_test.go:sidekiq-secret:126",
		"525d9792b1e3670b4630b8fcc385ca22e8544f9b:detect/detect_test.go:sidekiq-secret:142",
	}

	if !reflect.DeepEqual(test[:20], want) || err != nil {
		t.Fatalf(`loadIgnoredFingerprints(gitleaks, ...) = %v, %v, want %v, nil`, test, err, want)
	}
}
