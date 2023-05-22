package project

import (
	"time"
)

type ID uint64

// Project represents project top data
type Project struct {
	ID         ID
	Alias      string
	Language   string
	FolderPath string
	CreatedAt  time.Time
	// Add git path for Hosted version
}

// TODO packages prioritizing based on business value
// type Package struct {
// 	ID        ID
// 	ProjectID ID
// 	Path      string
// 	Priority  Priority
// 	Present   bool
// }
//
// type Priority string
//
// const (
// 	Regular    Priority = ""
// 	Vital               = "vital"
// 	Money               = "money"
// 	Critical            = "critical"
// 	Deprecated          = "deprecated"
// )

// TODO next
// Версия кода. При загрузке новых данных более свежего коммита, старые данные помечаются неактуальными
// type Revision struct {
// 	Hash string
// }

type File struct {
	ID      ID
	Project ID
	Package string
	Name    string
	Lines   uint32
	Symbols uint32
	Present bool
}

type GitCommit struct {
	ID      ID
	Hash    string
	Author  string
	Message string
	Time    time.Time
}

type GitChange struct {
	ID          ID
	File        ID
	Commit      ID
	RowsAdded   uint32
	RowsRemoved uint32
	Time        time.Time `gorm:"index:,sort:desc`
}

type Coverage struct {
	File           ID
	Percent        uint8
	UncoveredCount uint32
	UncoveredLines []uint32 `gorm:"serializer:json"`
}

// TODO future UI
// DataFetchJob contains project data collection job state
// type DataFetchJob struct {
// 	ProjectID  ID
// 	DataSource string
// 	CreatedAt  time.Time
// 	FinishedAt *time.Time
// }
