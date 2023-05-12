package git_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	git2 "devex_dashboard/datasource/git"
)

func createTestRepository(t *testing.T, commits, authors, files int) string {
	temp, err := os.MkdirTemp("", "TestExtract-*")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(temp))
	})

	repository, err := git.PlainInit(temp, false)
	require.NoError(t, err)

	worktree, err := repository.Worktree()
	require.NoError(t, err)

	for i := 0; i < commits; i++ {
		name := fmt.Sprintf("%d.txt", i%files)
		dir := strconv.Itoa(i % files)
		fullDir := filepath.Join(temp, dir)
		file := filepath.Join(fullDir, name)
		require.NoError(t, os.MkdirAll(fullDir, os.ModePerm))
		require.NoError(t, os.WriteFile(file, []byte(fmt.Sprintf("Test Content %d", i)), os.ModePerm))

		_, err = worktree.Add(filepath.Join(dir, name))
		require.NoError(t, err)

		_, err = worktree.Commit(fmt.Sprintf("%d commit", i), &git.CommitOptions{
			Author: &object.Signature{
				Name:  fmt.Sprintf("%d Name", i%authors),
				Email: fmt.Sprintf("%d@test.com", i%authors),
				When:  time.Now(),
			},
		})
		require.NoError(t, err)
	}

	return temp
}

func TestExtract(t *testing.T) {
	path := createTestRepository(t, 10, 3, 4)

	c := make(chan git2.FileCommit, 20)
	var (
		resultCommits int
		authors       = map[string]struct{}{}
		files         = map[string]struct{}{}
	)

	require.NoError(t, git2.ExtractCommits(context.TODO(), path, c))

	for commit := range c {
		resultCommits++
		authors[commit.Author] = struct{}{}
		files[commit.File] = struct{}{}
	}

	assert.Equal(t, 10, resultCommits)
	assert.Equal(t, 3, len(authors))
	assert.Equal(t, 4, len(files))
}
