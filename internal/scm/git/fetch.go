package git

import (
	"errors"
	"fmt"
	"math"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func Fetch(repo string, shallow bool) error {
	r, err := git.PlainOpen(repo)
	if err != nil {
		return err
	}

	opts := &git.FetchOptions{}
	if shallow {
		opts.RefSpecs = []config.RefSpec{"+HEAD:refs/remotes/origin/HEAD"}
		opts.Tags = git.NoTags
		opts.Depth = 1
	} else {
		opts.RefSpecs = []config.RefSpec{
			config.RefSpec(fmt.Sprintf(config.DefaultFetchRefSpec, "origin")),
		}
		opts.Tags = git.AllTags
		opts.Depth = math.MaxInt32 // means infinite depth, cf. https://git-scm.com/docs/shallow
	}

	err = r.Fetch(opts)
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return nil
		}
		return err
	}

	return nil
}
