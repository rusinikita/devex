package datacollector

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/rusinikita/devex/datasource"
	"github.com/rusinikita/devex/datasource/files"
	"github.com/rusinikita/devex/datasource/git"
	"github.com/rusinikita/devex/datasource/testcoverage"
	"github.com/rusinikita/devex/project"
)

// DepWheel chart
// go list -json="ImportPath,Imports" ./...

func Collect(ctx context.Context, db *gorm.DB, pkt project.Project, extractors datasource.Extractors) error {
	group, _ := errgroup.WithContext(ctx)

	log.Println("Start files data collection")

	fileChan := make(chan files.File)
	group.Go(func() error {
		for file := range fileChan {
			err := db.Create(&project.File{
				Package: file.Package,
				Name:    file.Name,
				Project: pkt.ID,
				Lines:   file.Lines,
				Symbols: file.Symbols,
				Tags:    file.Tags,
				Imports: file.Imports,
				Present: true,
			}).Error
			if err != nil {
				return fmt.Errorf("file saving: %q", err)
			}
		}

		return nil
	})

	err := extractors.Files(ctx, pkt.FolderPath, fileChan)
	if err != nil {
		return fmt.Errorf("files collection: %q", err)
	}

	if extractors.Coverage != nil {
		c := make(chan testcoverage.Package)

		group.Go(func() error {
			for pkg := range c {
				for _, file := range pkg.Files {
					projectFile := project.File{
						Name:    file.File,
						Package: pkg.Path,
						Project: pkt.ID,
					}
					tx := db.FirstOrCreate(&projectFile, projectFile)
					err := tx.Error
					if err != nil {
						return fmt.Errorf("finding coverage file: %q", err)
					}

					err = db.Create(&project.Coverage{
						File:           projectFile.ID,
						Percent:        file.Percent,
						UncoveredCount: uint32(len(file.UncoveredLines)),
						UncoveredLines: file.UncoveredLines,
					}).Error
					if err != nil {
						return fmt.Errorf("commit saving: %q", err)
					}
				}
			}

			return nil
		})

		log.Println("Start coverage data collection")

		err := extractors.Coverage(ctx, pkt.FolderPath, c)
		if err != nil {
			log.Printf("skip coverage collection: %q\n", err)
		}
	}

	if extractors.Git != nil {
		c := make(chan git.Commit)

		group.Go(func() error {
			defer close(c)
			commitsHandled := 0

			for commit := range c {
				commitsHandled++
				if commitsHandled%100 == 0 {
					log.Println(commitsHandled, "commits handled")
				}

				c := project.GitCommit{
					Hash:    commit.Hash,
					Author:  commit.Author,
					Message: commit.Message,
					Time:    commit.Time,
				}
				err := db.FirstOrCreate(&c, c).Error
				if err != nil {
					return fmt.Errorf("finding commit: %q", err)
				}

				for _, cFile := range commit.Files {
					file := project.File{
						Name:    cFile.File,
						Package: cFile.Package,
						Project: pkt.ID,
					}
					err := db.FirstOrCreate(&file, file).Error
					if err != nil {
						return fmt.Errorf("finding commit file: %q", err)
					}

					err = db.Create(&project.GitChange{
						File:        file.ID,
						Commit:      c.ID,
						RowsAdded:   cFile.RowsAdded,
						RowsRemoved: cFile.RowsRemoved,
						Time:        commit.Time,
					}).Error
					if err != nil {
						return fmt.Errorf("commit saving: %q", err)
					}
				}
			}

			return nil
		})

		log.Println("Start git data collection")

		err := extractors.Git(ctx, pkt.FolderPath, c)
		if err != nil {
			return fmt.Errorf("git commits collection: %q", err)
		}
	}

	log.Println("Done. Finishing")

	return group.Wait()
}
