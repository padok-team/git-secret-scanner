package scm

import (
	"context"
	"fmt"

	"github.com/padok-team/git-secret-scanner/internal/scm/git"
	"github.com/xanzy/go-gitlab"
)

type Visibility int

// repository visiblity
const (
	VisibilityAll Visibility = iota
	VisibilityPrivate
	VisibilityPublic
)

func (v Visibility) String() string {
	switch v {
	case VisibilityAll:
		return "all"
	case VisibilityPrivate:
		return "private"
	case VisibilityPublic:
		return "public"
	default:
		// should never be reached
		return ""
	}
}

// Visiblity must implement cobra pflag.Value interface
func (v *Visibility) Set(s string) error {
	switch s {
	case VisibilityAll.String():
		*v = VisibilityAll
		return nil
	case VisibilityPrivate.String():
		*v = VisibilityPrivate
		return nil
	case VisibilityPublic.String():
		*v = VisibilityPublic
		return nil
	default:
		return fmt.Errorf(
			"visiblity must be one of \"%s\", \"%s\" or \"%s\"",
			VisibilityAll.String(),
			VisibilityPrivate.String(),
			VisibilityPublic.String(),
		)
	}
}

// Visiblity must implement cobra pflag.Value interface
func (v *Visibility) Type() string {
	return "{all,private,public}"
}

func (v Visibility) Gitlab() *gitlab.VisibilityValue {
	var visiblity gitlab.VisibilityValue

	switch v {
	case VisibilityAll:
		return nil
	case VisibilityPublic:
		visiblity = gitlab.PublicVisibility
		return &visiblity
	case VisibilityPrivate:
		visiblity = gitlab.PrivateVisibility
		return &visiblity
	default:
		// should not be reached
		return nil
	}
}

type Scm interface {
	GitRepoUrl(repository string) string
	ListRepos(ctx context.Context) ([]string, error)
}

type ScmConfig struct {
	Org             string
	Server          string
	GitProtocol     git.Protocol
	Visiblity       Visibility
	IncludeArchived bool
}
