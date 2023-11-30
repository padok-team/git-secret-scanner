package git

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
)

func Clone(ctx context.Context, url string, destination string, shallow bool, fetchIfExists bool) error {
	opts := &git.CloneOptions{URL: url}
	if shallow {
		opts.SingleBranch = true
		opts.Tags = git.NoTags
		opts.Depth = 1
	}

	_, err := git.PlainCloneContext(ctx, destination, true, opts)
	if err != nil {
		// if the path already exists, it means the repostiory has already been cloned and it is not an error
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			if fetchIfExists {
				if err := Fetch(destination, shallow); err != nil {
					return fmt.Errorf("fetch error: %w", err)
				}
				return nil
			}
			return nil
		} else {
			return err
		}
	}

	return nil
}
