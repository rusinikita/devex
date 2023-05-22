package datasource

import (
	"context"

	"devex_dashboard/datasource/files"
	"devex_dashboard/datasource/git"
	"devex_dashboard/datasource/testcoverage"
	"devex_dashboard/project"
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
