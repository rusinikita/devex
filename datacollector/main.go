package datacollector

import (
	"context"
	"errors"
	"fmt"
	"github.com/rusinikita/devex/dao"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/rusinikita/devex/datasource"
	"github.com/rusinikita/devex/datasource/files"
	"github.com/rusinikita/devex/datasource/git"
	"github.com/rusinikita/devex/datasource/lint"
	"github.com/rusinikita/devex/datasource/testcoverage"
	"github.com/rusinikita/devex/project"
)

// DepWheel chart
// go list -json="ImportPath,Imports" ./...

type ErrStruct struct {
	Error    error
	Template string
}

const commitHandleCount = 100

func Collect( //nolint
	ctx context.Context,
	db *gorm.DB,
	pkt dao.Project,
	extractors datasource.Extractors,
) error {
	log.Println("Start files data collection")

	group, _ := errgroup.WithContext(ctx)

	fileChan := make(chan files.File)
	group.Go(groupFuncFile(db, pkt.ID, fileChan))

	err := extractors.Files(ctx, pkt.FolderPath, fileChan)
	if err != nil {
		return fmt.Errorf("files collection: %q", err)
	}

	if extractors.Coverage != nil {
		c := make(chan testcoverage.Package)

		group.Go(groupFuncCoverage(db, pkt.ID, c))

		log.Println("Start coverage data collection")

		err := extractors.Coverage(ctx, pkt.FolderPath, c)
		if err != nil {
			log.Printf("skip coverage collection: %q\n", err)
		}
	}

	if extractors.Git != nil {
		c := make(chan git.Commit)

		group.Go(groupFuncCommit(db, pkt.ID, c))

		log.Println("Start git data collection")

		err := extractors.Git(ctx, pkt.FolderPath, c)
		if err != nil {
			return fmt.Errorf("git commits collection: %q", err)
		}
	}

	log.Println("Done. Finishing")

	return group.Wait()
}

func groupFuncFile(db *gorm.DB, projectID project.ID, fileChan chan files.File) func() error {
	return func() error {
		for file := range fileChan {
			err := saveFile(db, projectID, file)
			if err != nil {
				return fmt.Errorf("file saving: %q", err)
			}
		}

		return nil
	}
}

func groupFuncCommit(db *gorm.DB, id project.ID, c chan git.Commit) func() error {
	return func() error {
		commitsHandled := 0

		for commit := range c {
			commitsHandled++

			err := processingSaveCommit(commitsHandled, db, commit, id)

			if err != nil {
				return err
			}
		}

		return nil
	}
}

func processingSaveCommit(commitsHandled int, db *gorm.DB, commit git.Commit, id project.ID) error {
	if commitsHandled%commitHandleCount == 0 {
		log.Println(commitsHandled, "commits handled")
	}
	gitCommit, err := saveGitCommit(db, commit)

	if err != nil {
		return fmt.Errorf("finding commit: %q", err)
	}

	for _, cFile := range commit.Files {
		errStructure := saveFileAndGitChange(db, cFile, id, gitCommit)

		if errStructure.Error != nil {
			return fmt.Errorf(errStructure.Template, errStructure.Error)
		}
	}
	return nil
}

func saveFileAndGitChange(db *gorm.DB, cFile git.FileCommit, id project.ID, commit dao.GitCommit) ErrStruct {
	file, errorStruct := saveFileByParams(db, id, cFile.Package, cFile.File)

	if errorStruct.Error != nil {
		return errorStruct
	}

	return saveGitChange(db, cFile, file, commit)
}

func saveGitChange(db *gorm.DB, cFile git.FileCommit, file dao.File, commit dao.GitCommit) ErrStruct {
	err := db.Create(&dao.GitChange{
		File:        file.ID,
		Commit:      commit.ID,
		RowsAdded:   cFile.RowsAdded,
		RowsRemoved: cFile.RowsRemoved,
		Time:        commit.Time,
	}).Error

	if err != nil {
		return ErrStruct{
			Template: "commit saving: %q",
			Error:    err,
		}
	}

	return ErrStruct{}
}

func saveGitCommit(db *gorm.DB, commit git.Commit) (dao.GitCommit, error) {
	gitCommit := dao.GitCommit{
		Hash:    commit.Hash,
		Author:  commit.Author,
		Message: commit.Message,
		Time:    commit.Time,
	}
	err := db.FirstOrCreate(&gitCommit, gitCommit).Error

	return gitCommit, err
}

func groupFuncCoverage(db *gorm.DB, projectID project.ID, c chan testcoverage.Package) func() error {
	return func() error {
		for pkg := range c {
			for _, file := range pkg.Files {
				err := saveFileAndCoverage(db, file, projectID, pkg.Path)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}

func saveFileAndCoverage(db *gorm.DB, file testcoverage.Coverage, id project.ID, path string) error {
	projectFile, errStructure := saveFileByParams(db, id, path, file.File)

	if errStructure.Error != nil {
		return fmt.Errorf(errStructure.Template, errStructure.Error)
	}

	errStructure = saveCoverage(db, file, projectFile.ID)

	if errStructure.Error != nil {
		return fmt.Errorf(errStructure.Template, errStructure.Error)
	}
	return nil
}

func saveCoverage(db *gorm.DB, file testcoverage.Coverage, id project.ID) ErrStruct {
	err := db.Create(&dao.Coverage{
		File:           id,
		Percent:        file.Percent,
		UncoveredCount: uint32(len(file.UncoveredLines)),
		UncoveredLines: file.UncoveredLines,
	}).Error
	if err != nil {
		return ErrStruct{Template: "commit saving: %q", Error: err}
	}
	return ErrStruct{}
}

func saveFileByParams(db *gorm.DB, id project.ID, path string, file string) (dao.File, ErrStruct) {
	projectFile := dao.File{
		Name:    file,
		Package: path,
		Project: id,
	}
	tx := db.FirstOrCreate(&projectFile, projectFile)
	err := tx.Error
	if err != nil {
		return projectFile, ErrStruct{
			Template: "finding coverage file: %q",
			Error:    err,
		}
	}

	return projectFile, ErrStruct{}
}

func saveFile(db *gorm.DB, projectID project.ID, file files.File) error {
	return db.Create(&dao.File{
		Package: file.Package,
		Name:    file.Name,
		Project: projectID,
		Lines:   file.Lines,
		Symbols: file.Symbols,
		Tags:    file.Tags,
		Imports: file.Imports,
		Present: true,
	}).Error
}

func CheckStyle(database *gorm.DB, projectAlias string, filePath string) error {
	projectID, err := getProjectIDByAlias(database, projectAlias)
	if err != nil {
		return err
	}

	log.Printf("project found, id: %d \n", projectID)

	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	return extractAndBatchRows(database, projectID, file)
}

func extractAndBatchRows(database *gorm.DB, projectID project.ID, file *os.File) error {
	lintFiles, err := lint.ExtractCheckStyleXML(file)

	if err != nil {
		return err
	}
	log.Printf("count files in report: %d \n", len(lintFiles))

	return batchRows(database, projectID, lintFiles)
}

func getProjectIDByAlias(database *gorm.DB, alias string) (project.ID, error) {
	var projectDao dao.Project
	tx := database.Select("id").Where("alias = ?", alias).Take(&projectDao)

	return projectDao.ID, tx.Error
}

func batchRows(database *gorm.DB, projectID project.ID, lintFiles []lint.LinterFile) error {
	return database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Error; err != nil {
			return err
		}

		var filesFromDB []dao.File
		tx.Where("project = ?", projectID).Find(&filesFromDB)

		fileIds := getFileIdsByDBFiles(filesFromDB)

		tx.Where("file_id IN ?", fileIds).Delete(&dao.LintError{})

		lintErrors, err := convertDtoToDao(lintFiles, filesFromDB)
		if err != nil {
			return err
		}

		return tx.Create(&lintErrors).Error
	})
}

func getFileIdsByDBFiles(filesFromDB []dao.File) []project.ID {
	var fileIds []project.ID
	for _, file := range filesFromDB {
		fileIds = append(fileIds, file.ID)
	}
	return fileIds
}

func convertDtoToDao(lintFiles []lint.LinterFile, fileDaoList []dao.File) ([]dao.LintError, error) {
	var result []dao.LintError

	for _, lintFile := range lintFiles {
		fileID, err := getFileIDByPath(fileDaoList, lintFile.Path)
		if err != nil {
			return result, err
		}

		for _, lintError := range lintFile.Errors {
			result = append(result, getLintError(fileID, lintError))
		}
	}

	return result, nil
}

func getLintError(fileID project.ID, lintError lint.LinterError) dao.LintError {
	return dao.LintError{
		FileID:     fileID,
		FileColumn: lintError.Column,
		FileLine:   lintError.Line,
		Message:    lintError.Message,
		Severity:   lintError.Severity,
		Source:     lintError.Source,
	}
}

func getFileIDByPath(projectFiles []dao.File, path string) (project.ID, error) {
	for _, file := range projectFiles {
		// todo fixed on the issue #9
		if filepath.Join(file.Package, file.Name) == path {
			return file.ID, nil
		}
		// todo fixed on the issue #9
	}

	message := fmt.Sprintf("internal error, file id not found for path %s", path)

	return 0, errors.New(message)
}
