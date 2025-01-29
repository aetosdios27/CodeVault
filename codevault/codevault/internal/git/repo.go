package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func InitRepo(path string) (*git.Repository, error) {
	return git.PlainInit(path, false)
}

func CommitFile(repo *git.Repository, path, message string) error {
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}
	_, err = w.Add(path)
	if err != nil {
		return fmt.Errorf("failed to add file for commit: %w", err)
	}
	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "CodeVault",
			Email: "auto@codevault",
		},
	})
	if err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}
