package dao

import "github.com/rusinikita/devex/project"

type Coverage struct {
	UncoveredLines []uint32 `gorm:"serializer:json"`
	Percent        uint8
	UncoveredCount uint32
	File           project.ID
}
