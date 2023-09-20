package dashboard

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rusinikita/devex/db"
	"github.com/rusinikita/devex/project"
)

func Test_FileSizes(t *testing.T) {
	database := db.TestDB("../devex.db")

	gotResult, err := FileSizes(database, []project.ID{1}, "and (name like '%.go' or name like '%.md')")
	require.NoError(t, err)

	log.Println(gotResult.BarNames())
}
