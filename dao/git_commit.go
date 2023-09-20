package dao

import (
	"time"

	"github.com/rusinikita/devex/project"
)

type GitCommit struct {
	Time    time.Time
	Hash    string
	Author  string
	Message string
	ID      project.ID
}
