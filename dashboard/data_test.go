package dashboard

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinikita/devex/db"
	"github.com/rusinikita/devex/project"
)

func TestName2(t *testing.T) {
	t.Skip("local test")

	result, err := commitMessages(db.TestDB("../devex.db").Debug(), false, []project.ID{1}, "", "")

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
