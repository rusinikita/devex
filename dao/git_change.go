package dao

import (
	"time"

	"github.com/rusinikita/devex/project"
)

type GitChange struct {
	Time        time.Time `gorm:"index:,sort:desc"`
	ID          project.ID
	File        project.ID
	Commit      project.ID
	RowsAdded   uint32
	RowsRemoved uint32
}
