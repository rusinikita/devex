package lint

import (
	"encoding/xml"
	"io"
)

type xmlFile struct {
	Name        xml.Name     `xml:"checkstyle"`
	LinterFiles []LinterFile `xml:"file"`
}

type LinterFile struct {
	Path   string        `xml:"name,attr"`
	Errors []LinterError `xml:"error"`
}
type LinterError struct {
	Column   uint   `xml:"column,attr"`
	Line     uint   `xml:"line,attr"`
	Message  string `xml:"message,attr"`
	Severity string `xml:"severity,attr"`
	Source   string `xml:"source,attr"`
}

// ExtractCheckStyleXml парсим xml
func ExtractCheckStyleXml(file io.Reader) ([]LinterFile, error) {
	var data xmlFile

	err := xml.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.LinterFiles, nil
}
