package lint

import (
	"encoding/xml"
	"io"
)

type xmlFile struct {
	CheckStyle xml.Name `xml:"checkstyle"`
}

type LinterFile struct {
	Path   string
	Errors []LinterError
}
type LinterError struct {
	Column   uint
	Line     uint
	Message  string
	Severity string
	Source   string
}

func extractCheckStyleXml(file io.Reader) ([]LinterFile, error) {
	var data xmlFile

	err := xml.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	return []LinterFile{}, nil
}
