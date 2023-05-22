package dashboard

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"devex_dashboard/db"
	"devex_dashboard/project"
)

func TestName2(t *testing.T) {
	t.Skip("local test")

	result, err := contribution(db.TestDB("../devex_bd.db").Debug(), []project.ID{1}, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestTreemap(t *testing.T) {
	v := values{
		valueData{
			Alias:   "a",
			Package: "b/c/d",
			Name:    "file.go",
			Value:   10,
		},
	}

	result := v.treeMaps()
	assert.NotEmpty(t, result)
	assert.NotEmpty(t, result[0].Children)
}
