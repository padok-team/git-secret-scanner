package git

import "fmt"

type Protocol int

// git protocol
const (
	ProtocolHTTPS Protocol = iota
	ProtocolSSH
)

func (p Protocol) String() string {
	switch p {
	case ProtocolHTTPS:
		return "https"
	case ProtocolSSH:
		return "ssh"
	default:
		// should never be reached
		return ""
	}
}

func RepoUrl(repository string, protocol Protocol, server string, token string) string {
	url := ""
	switch protocol {
	case ProtocolHTTPS:
		if token == "" {
			url = fmt.Sprintf("https://%s/%s", server, repository)
		} else {
			url = fmt.Sprintf("https://x-access-token:%s@%s/%s", token, server, repository)
		}
	case ProtocolSSH:
		url = fmt.Sprintf("git@%s:%s.git", server, repository)
	default:
	}

	return url
}
