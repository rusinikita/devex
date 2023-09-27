package view

import (
	"github.com/rusinikita/devex/dao"
	"github.com/rusinikita/devex/project"
	"github.com/rusinikita/devex/slices"
)

type formStruct struct {
	SelectedProjects slices.Set[project.ID]
	PackageFilter    string
	NameFilter       string
	TrimPackage      string
	CommitFilters    string
	FileFilters      string
	Projects         []dao.Project
	PerFiles         bool
	PerFilesImports  bool
}

func newFormStruct(projects []dao.Project, params Params) formStruct {
	return formStruct{
		Projects:         projects,
		SelectedProjects: slices.ToSet(params.ProjectIDs),
		PerFiles:         params.PerFiles,
		PerFilesImports:  params.PerFilesImports,
		PackageFilter:    params.PackageFilter,
		NameFilter:       params.NameFilter,
		TrimPackage:      params.TrimPackage,
		CommitFilters:    params.CommitFilters,
		FileFilters:      params.FileFilters,
	}
}
