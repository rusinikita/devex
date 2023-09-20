package dao

import "github.com/rusinikita/devex/project"

type File struct {
	Tags    map[string]uint32 `gorm:"serializer:json"`
	Package string
	Name    string
	Imports []string `gorm:"serializer:json"`
	Present bool
	Lines   uint32
	Symbols uint32
	ID      project.ID
	Project project.ID
}
