package datacollector_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"devex_dashboard/datacollector"
	"devex_dashboard/datasource"
	"devex_dashboard/datasource/files"
	"devex_dashboard/datasource/git"
	"devex_dashboard/datasource/testcoverage"
	"devex_dashboard/db"
	"devex_dashboard/project"
)

func TestCollect(t *testing.T) {
	ctx := context.TODO()
	m := mock.Mock{}
	m.On("Files", mock.Anything).Return(nil)
	m.On("Git", mock.Anything).Return(nil)
	m.On("Coverage", mock.Anything).Return(nil)

	e := datasource.Extractors{
		Files: func(ctx context.Context, projectPath string, c chan<- files.File) error {
			defer close(c)
			m.MethodCalled("Files", projectPath)

			for i := 0; i < 10; i++ {
				c <- files.File{
					Package: strconv.Itoa(i % 3),
					Name:    strconv.Itoa(i),
					Lines:   10,
					Symbols: 100,
				}
			}

			return nil
		},
		Git: func(ctx context.Context, projectPath string, c chan<- git.FileCommit) error {
			defer close(c)
			m.MethodCalled("Git", projectPath)

			for i := 0; i < 30; i++ {
				c <- git.FileCommit{
					Package:     strconv.Itoa(i % 3),
					File:        strconv.Itoa(i / 3),
					Hash:        strconv.Itoa(i),
					Author:      "test",
					RowsAdded:   1,
					RowsRemoved: 1,
					Time:        time.Now(),
				}
			}

			return nil
		},
		Coverage: func(ctx context.Context, projectPath string, c chan<- testcoverage.Package) error {
			defer close(c)
			m.MethodCalled("Coverage", projectPath)

			for i := 0; i < 3; i++ {
				p := testcoverage.Package{
					Path:  strconv.Itoa(i),
					Files: nil,
				}

				for i := 0; i < 10; i++ {
					p.Files = append(p.Files, testcoverage.Coverage{
						File:           strconv.Itoa(i),
						Percent:        uint8(i) * 10,
						UncoveredLines: []uint32{1, 2, 3, 4, 5},
					})
				}

				c <- p
			}

			return nil
		},
	}

	database := db.TestDB()

	p := project.Project{
		ID:         1,
		Alias:      "test",
		Language:   "go",
		FolderPath: "test/test",
		CreatedAt:  time.Time{},
	}

	err := datacollector.Collect(ctx, database, p, e)
	require.NoError(t, err)

	m.AssertCalled(t, "Git", p.FolderPath)
	m.AssertCalled(t, "Files", p.FolderPath)
	m.AssertCalled(t, "Coverage", p.FolderPath)

	var resultFiles []project.File

	assert.NoError(t, database.Find(&resultFiles, project.File{Present: true}).Error)
	assert.Len(t, resultFiles, 10)

	var commits []project.GitChange

	assert.NoError(t, database.Find(&commits).Error)
	assert.Len(t, commits, 30)

	var coverages []project.Coverage

	assert.NoError(t, database.Find(&coverages).Error)
	assert.Len(t, coverages, 30)
	assert.Equal(t, []uint32{1, 2, 3, 4, 5}, coverages[0].UncoveredLines)

	resultFiles = nil
	assert.NoError(t, database.Find(&resultFiles).Error)
	assert.Len(t, resultFiles, 30)
}