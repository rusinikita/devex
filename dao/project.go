package dao

import (
	"time"

	"github.com/rusinikita/devex/project"
)

// Project represents project top data
type Project struct {
	CreatedAt  time.Time
	Alias      string
	Language   string
	FolderPath string
	ID         project.ID
	// Add git path for Hosted version
}
