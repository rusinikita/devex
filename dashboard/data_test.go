package dashboard

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinikita/devex/db"
	"github.com/rusinikita/devex/project"
)

func TestName2(t *testing.T) {
	t.Skip("local test")

	result, err := CommitMessages(db.TestDB("../devex.db").Debug(), false, []project.ID{1}, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// todo failure test
// Error:      	Should NOT be empty, but was []
func TestGetResult(t *testing.T) {
	listValues := Values{
		ValueData{Value: 3},
		ValueData{Value: 1},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
	}

	result := GetResult(listValues)

	assert.Len(t, result, 20)
	assert.Equal(t, result, Values{
		ValueData{Value: 3},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
		ValueData{Value: 2},
	})
}
