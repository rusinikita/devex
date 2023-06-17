package dashboard

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"devex_dashboard/db"
	"devex_dashboard/project"
)

func Test_fileSizes(t *testing.T) {
	database := db.TestDB("../devex.db")

	gotResult, err := fileSizes(database, []project.ID{1}, "and (name like '%.go' or name like '%.md')")
	require.NoError(t, err)

	fmt.Println(gotResult.barNames())
}
