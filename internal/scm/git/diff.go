package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const ErrCommitDiff string = "failed to get commit diff: %w"

func CommitDiff(repo string, hash string, file string) (string, error) {
	r, err := git.PlainOpen(repo)
	if err != nil {
		return "", fmt.Errorf(ErrCommitDiff, err)
	}

	commit, err := r.CommitObject(plumbing.NewHash(hash))
	if err != nil {
		return "", fmt.Errorf(ErrCommitDiff, err)
	}

	f, err := commit.File(file)
	if err != nil {
		return "", fmt.Errorf(ErrCommitDiff, err)
	}

	diff, err := f.Contents()
	if err != nil {
		return "", fmt.Errorf(ErrCommitDiff, err)
	}

	return diff, nil
}
