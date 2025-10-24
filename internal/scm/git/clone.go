package git

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
)

func Clone(ctx context.Context, url string, destination string, shallow bool, bare bool, fetchIfExists bool) error {
	opts := &git.CloneOptions{URL: url}
	if shallow {
		opts.SingleBranch = true
		opts.Tags = git.NoTags
		opts.Depth = 1
	}

	_, err := git.PlainCloneContext(ctx, destination, bare, opts)
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
	} else if !shallow {
		// if it's not a shallow clone, fetch all other remote refs that are not retrieved by default
		if err := Fetch(destination, false); err != nil {
			return fmt.Errorf("fetch error: %w", err)
		}
	}

	return nil
}
