package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const ErrReadFileContent string = "failed to read file content: %w"

func readCommitFile(r *git.Repository, hash plumbing.Hash, file string) (string, error) {
	commit, err := r.CommitObject(hash)
	if err != nil {
		return "", fmt.Errorf(ErrReadFileContent, err)
	}

	f, err := commit.File(file)
	if err != nil {
		return "", fmt.Errorf(ErrReadFileContent, err)
	}

	diff, err := f.Contents()
	if err != nil {
		return "", fmt.Errorf(ErrReadFileContent, err)
	}

	return diff, nil
}

func ReadCommitFileContent(repo string, hash string, file string) (string, error) {
	r, err := git.PlainOpen(repo)
	if err != nil {
		return "", fmt.Errorf(ErrReadFileContent, err)
	}

	return readCommitFile(r, plumbing.NewHash(hash), file)
}

func ReadFileContent(repo string, file string) (string, error) {
	r, err := git.PlainOpen(repo)
	if err != nil {
		return "", fmt.Errorf(ErrReadFileContent, err)
	}

	head, err := r.Head()
	if err != nil {
		return "", fmt.Errorf(ErrReadFileContent, err)
	}

	return readCommitFile(r, head.Hash(), file)
}
