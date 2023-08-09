package datasource

import (
	"context"

	"github.com/rusinikita/devex/datasource/files"
	"github.com/rusinikita/devex/datasource/git"
	"github.com/rusinikita/devex/datasource/testcoverage"
	"github.com/rusinikita/devex/project"
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
		Coverage: testcoverage.ExtractXml,
	}
}

func DataEntities() []any {
	return []any{
		project.Project{},
		project.File{},
		project.Coverage{},
		project.GitChange{},
		project.GitCommit{},
	}
}
