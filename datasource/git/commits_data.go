package git

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Commit struct {
	Hash    string
	Author  string
	Message string
	Files   []FileCommit
	Time    time.Time
}

type FileCommit struct {
	Package     string
	File        string
	RowsAdded   uint32
	RowsRemoved uint32
}

func ExtractCommits(ctx context.Context, projectPath string, c chan<- Commit) error {
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

		var files []FileCommit
		for _, file := range stats {
			files = append(files, FileCommit{
				Package:     strings.TrimPrefix(filepath.Dir(file.Name), "."),
				File:        filepath.Base(file.Name),
				RowsAdded:   uint32(file.Addition),
				RowsRemoved: uint32(file.Deletion),
			})
		}

		c <- Commit{
			Hash:    commit.Hash.String(),
			Author:  commit.Author.Email,
			Message: commit.Message,
			Time:    commit.Author.When,
			Files:   files,
		}

		return nil
	})
}
