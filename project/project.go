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
	Tags    map[string]uint32 `gorm:"serializer:json"` // experiment with tags: nolint,billing,money,order
	Imports []string          `gorm:"serializer:json"`
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
	Time        time.Time `gorm:"index:,sort:desc"`
}

type Import struct {
	File ID
	Path string
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

type LintError struct {
	Id         ID        `gorm:"primaryKey"`
	FileId     ID        `gorm:"column:file_id;not null;index;comment:Foreign key to files"`
	CreatedAt  time.Time `gorm:"column:created_at;default:(DATETIME('now'));not null;comment:created at"`
	FileColumn uint      `gorm:"column:file_column;not null;comment:Column with error"`
	FileLine   uint      `gorm:"column:file_line;not null;comment:Row with error"`
	Message    string    `gorm:"column:message;type:text;not null;comment:Error message"`
	Severity   string    `gorm:"column:message;type:varchar(155);not null;comment:Severity error"`
	Source     string    `gorm:"column:message;type:varchar(155);not null;comment:What source found error"`
	File       *File     `gorm:"foreignKey:FileId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
