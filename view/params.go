package view

import (
	"github.com/rusinikita/devex/project"
	"github.com/rusinikita/devex/slices"
)

// Params TODO it is request?
type Params struct {
	PackageFilter   string       `form:"package_filter"`
	NameFilter      string       `form:"name_filter"`
	TrimPackage     string       `form:"trim_package"`
	CommitFilters   string       `form:"commit_filters"`
	FileFilters     string       `form:"file_filters"`
	ProjectIDs      []project.ID `form:"project_ids"`
	PerFiles        bool         `form:"per_files"`
	PerFilesImports bool         `form:"per_files_imports"`
}

func (p Params) sqlFilter() (sql string) {
	if p.PackageFilter != "" {
		sql += "and " + slices.SQLFilter("package", p.PackageFilter)
	}

	if p.NameFilter != "" {
		sql += " and " + slices.SQLFilter("name", p.NameFilter)
	}

	return sql
}
