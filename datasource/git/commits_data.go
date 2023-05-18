package git

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type FileCommit struct {
	Package     string
	File        string
	Hash        string
	Author      string
	RowsAdded   uint32
	RowsRemoved uint32
	Time        time.Time
}

func ExtractCommits(ctx context.Context, projectPath string, c chan<- FileCommit) error {
	defer close(c)

	repository, err := git.PlainOpen(projectPath)
	if err != nil {
		return err
	}

	commitObjects, err := repository.CommitObjects()
	if err != nil {
		return err
	}

	return commitObjects.ForEach(func(commit *object.Commit) error {
		select {
		case <-ctx.Done():
			commitObjects.Close()
			return nil
		default:
		}

		stats, err := commit.Stats()
		if err != nil {
			return err
		}

		for _, file := range stats {
			raw := FileCommit{
				Package:     strings.TrimPrefix(filepath.Dir(file.Name), "."),
				File:        filepath.Base(file.Name),
				Hash:        commit.Hash.String(),
				Author:      commit.Author.Email,
				RowsAdded:   uint32(file.Addition),
				RowsRemoved: uint32(file.Deletion),
				Time:        commit.Author.When,
			}

			c <- raw
		}

		return nil
	})
}
