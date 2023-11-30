package git

import "testing"

func TestProtocolString(t *testing.T) {
	test := ProtocolHTTPS.String()
	want := "https"

	if test != want {
		t.Fatalf(`ProtocolHTTPS.String() = %q, want %q`, test, want)
	}

	test = ProtocolSSH.String()
	want = "ssh"

	if test != want {
		t.Fatalf(`ProtocolSSH.String() = %q, want %q`, test, want)
	}
}

func TestRepoUrl(t *testing.T) {
	test := RepoUrl("test_user/test_repo", ProtocolHTTPS, "test.server", "faketoken")
	want := "https://x-access-token:faketoken@test.server/test_user/test_repo"

	if test != want {
		t.Fatalf(`RepoUrl("test_user/test_repo", ProtocolHTTPS, "test.server", "faketoken") = %q, want %q`, test, want)
	}

	test = RepoUrl("test_user/test_repo", ProtocolSSH, "test.server", "faketoken")
	want = "git@test.server:test_user/test_repo.git"

	if test != want {
		t.Fatalf(`RepoUrl("test_user/test_repo", ProtocolSSH, "test.server", "faketoken") = %q, want %q`, test, want)
	}
}
