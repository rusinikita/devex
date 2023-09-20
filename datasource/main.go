package datasource

import (
	"context"

	"github.com/rusinikita/devex/dao"
	"github.com/rusinikita/devex/datasource/files"
	"github.com/rusinikita/devex/datasource/git"
	"github.com/rusinikita/devex/datasource/testcoverage"
)

type Extractor[T any] func(ctx context.Context, projectPath string, c chan<- T) error

type Extractors struct {
	Files    Extractor[files.File]
	Git      Extractor[git.Commit]
	Coverage Extractor[testcoverage.Package]
}

func NewExtractors() Extractors {
	return Extractors{
		Files:    files.Extract,
		Git:      git.ExtractCommits,
		Coverage: testcoverage.ExtractXMLCommand,
	}
}

func DataEntities() []any {
	return []any{
		dao.Project{},
		dao.File{},
		dao.Coverage{},
		dao.GitChange{},
		dao.GitCommit{},
		dao.LintError{},
	}
}
