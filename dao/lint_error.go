package dao

import (
	"time"

	"github.com/rusinikita/devex/project"
)

type LintError struct {
	File       *File      `gorm:"foreignKey:FileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  time.Time  `gorm:"column:created_at;default:(DATETIME('now'));not null;comment:created at"`
	Message    string     `gorm:"column:message;type:text;not null;comment:Error message"`
	Severity   string     `gorm:"column:message;type:varchar(155);not null;comment:Severity error"`
	Source     string     `gorm:"column:message;type:varchar(155);not null;comment:What source found error"`
	ID         project.ID `gorm:"primaryKey"`
	FileID     project.ID `gorm:"column:file_id;not null;index;comment:Foreign key to files"`
	FileColumn uint       `gorm:"column:file_column;not null;comment:Column with error"`
	FileLine   uint       `gorm:"column:file_line;not null;comment:Row with error"`
}
