package lint

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	one  = 1
	two  = 2
	four = 4
	nine = 9
)

func TestExtractCheckStyleXml(t *testing.T) {
	expect := []LinterFile{
		{
			Path: "slices/slices_test.go",
			Errors: []LinterError{
				{
					Column:   one,
					Line:     nine,
					Message:  "Function TestSQLFilter missing the call to method parallel\n",
					Severity: "error",
					Source:   "paralleltest",
				},
				{
					Column:   nine,
					Line:     one,
					Message:  "package should be 'slices_test' instead of 'slices'",
					Severity: "error",
					Source:   "testpackage",
				},
			},
		},
		{
			Path: "db/create.go",
			Errors: []LinterError{
				{
					Column:   two,
					Line:     four,
					Message:  "import 'github.com/glebarez/sqlite' is not allowed from list 'Main'",
					Severity: "error",
					Source:   "depguard",
				},
			},
		},
	}

	result, err := ExtractCheckStyleXML(bytes.NewBufferString(testFile))
	assert.NoError(t, err)

	assert.Equal(t, expect, result)
}

const testFile = `
	<?xml version="1.0" encoding="UTF-8"?>
	
	<checkstyle version="5.0">
  		<file name="slices/slices_test.go">
  		  <error column="1" line="9" message="Function TestSQLFilter missing the call to method parallel&#xA;" severity="error" source="paralleltest"></error>
  		  <error column="9" line="1" message="package should be 'slices_test' instead of 'slices'" severity="error" source="testpackage"></error>
  		</file>
  		<file name="db/create.go">
  		  <error column="2" line="4" message="import &#39;github.com/glebarez/sqlite&#39; is not allowed from list &#39;Main&#39;" severity="error" source="depguard"></error>
  		</file>
	</checkstyle>
`
