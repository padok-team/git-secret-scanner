package scm

import (
	"reflect"
	"testing"

	"github.com/padok-team/git-secret-scanner/internal/scm/git"
)

func TestNewGithubScm(t *testing.T) {
	config := &ScmConfig{
		Org:             "test_org",
		Server:          "",
		GitProtocol:     git.ProtocolHTTPS,
		Visiblity:       VisibilityAll,
		IncludeArchived: true,
	}

	test, err := NewGithubScm(config, "faketoken")

	if !reflect.DeepEqual(test.ScmConfig, config) || test.token != "faketoken" || test.client.BaseURL.String() != "https://api.github.com/" || err != nil {
		t.Fatalf(`NewGithubScm(...) = %v, %v, want {%v faketoken client}, nil`, test, err, config)
	}
	if config.Server != "github.com" {
		t.Fatalf(`NewGithubScm(...) -> config.Server = %q, want "github.com"`, config.Server)
	}

	config.Server = "test.server"
	test, err = NewGithubScm(config, "faketoken")

	if !reflect.DeepEqual(test.ScmConfig, config) || test.token != "faketoken" || test.client.BaseURL.String() != "https://test.server/api/v3/" || err != nil {
		t.Fatalf(`NewGithubScm(...) = %v, %v, want {%v faketoken client}, nil`, test, err, config)
	}
	if config.Server != "test.server" {
		t.Fatalf(`NewGithubScm(...) -> config.Server = %q, want "test.server"`, config.Server)
	}
}

func TestGithubScmGitRepoUrl(t *testing.T) {
	config := &ScmConfig{
		Org:             "test_org",
		Server:          "test.server",
		GitProtocol:     git.ProtocolHTTPS,
		Visiblity:       VisibilityAll,
		IncludeArchived: true,
	}

	scm, err := NewGithubScm(config, "faketoken")
	if err != nil {
		t.Fatalf(`NewGithubScm(...) = _, %v, want _, nil`, err)
	}

	test := scm.GitRepoUrl("test_user/test_repo")
	want := "https://x-access-token:faketoken@test.server/test_user/test_repo"

	if test != want {
		t.Fatalf(`scm.GitRepoUrl("test_user/test_repo") = %q, want %q`, test, want)
	}

	config.Server = "github.com"
	config.GitProtocol = git.ProtocolSSH

	scm, err = NewGithubScm(config, "faketoken")
	if err != nil {
		t.Fatalf(`NewGithubScm(...) = _, %v, want _, nil`, err)
	}

	test = scm.GitRepoUrl("test_user/test_repo")
	want = "git@github.com:test_user/test_repo.git"

	if test != want {
		t.Fatalf(`scm.GitRepoUrl("test_user/test_repo") = %q, want %q`, test, want)
	}
}
