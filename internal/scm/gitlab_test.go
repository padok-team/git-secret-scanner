package scm

import (
	"reflect"
	"testing"

	"github.com/padok-team/git-secret-scanner/internal/scm/git"
)

func TestNewGitlabScm(t *testing.T) {
	config := &ScmConfig{
		Org:             "test_org",
		Server:          "",
		GitProtocol:     git.ProtocolHTTPS,
		Visiblity:       VisibilityAll,
		IncludeArchived: true,
	}

	test, err := NewGitlabScm(config, "faketoken")

	if !reflect.DeepEqual(test.ScmConfig, config) || test.token != "faketoken" || test.client.BaseURL().String() != "https://gitlab.com/api/v4/" || err != nil {
		t.Fatalf(`NewGitlabScm(...) = %v, %v, want {%v faketoken client}, nil`, test, err, config)
	}
	if config.Server != "gitlab.com" {
		t.Fatalf(`NewGitlabScm(...) -> config.Server = %q, want "gitlab.com"`, config.Server)
	}

	config.Server = "test.server"
	test, err = NewGitlabScm(config, "faketoken")

	if !reflect.DeepEqual(test.ScmConfig, config) || test.token != "faketoken" || test.client.BaseURL().String() != "https://test.server/api/v4/" || err != nil {
		t.Fatalf(`NewGitlabScm(...) = %v, %v, want {%v faketoken client}, nil`, test, err, config)
	}
	if config.Server != "test.server" {
		t.Fatalf(`NewGitlabScm(...) -> config.Server = %q, want "test.server"`, config.Server)
	}
}

func TestGitlabScmGitRepoUrl(t *testing.T) {
	config := &ScmConfig{
		Org:             "test_org",
		Server:          "test.server",
		GitProtocol:     git.ProtocolHTTPS,
		Visiblity:       VisibilityAll,
		IncludeArchived: true,
	}

	scm, err := NewGitlabScm(config, "faketoken")
	if err != nil {
		t.Fatalf(`NewGitlabScm(...) = _, %v, want _, nil`, err)
	}

	test := scm.GitRepoUrl("test_user/test_repo")
	want := "https://x-access-token:faketoken@test.server/test_user/test_repo"

	if test != want {
		t.Fatalf(`scm.GitRepoUrl("test_user/test_repo") = %q, want %q`, test, want)
	}

	config.Server = "gitlab.com"
	config.GitProtocol = git.ProtocolSSH

	scm, err = NewGitlabScm(config, "faketoken")
	if err != nil {
		t.Fatalf(`NewGitlabScm(...) = _, %v, want _, nil`, err)
	}

	test = scm.GitRepoUrl("test_user/test_repo")
	want = "git@gitlab.com:test_user/test_repo.git"

	if test != want {
		t.Fatalf(`scm.GitRepoUrl("test_user/test_repo") = %q, want %q`, test, want)
	}
}
