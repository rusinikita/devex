package dashboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLFilter(t *testing.T) {
	expected := "(name like '%bla%' or name not like '%opa%') and name like '%dich%'"
	filter := "bla,!opa;dich"

	assert.Equal(t, expected, SQLFilter("name", filter))
}
