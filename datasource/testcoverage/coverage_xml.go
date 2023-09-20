package testcoverage

import (
	"context"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	Path  string
	Files []Coverage
}

type Coverage struct {
	File           string
	UncoveredLines []uint32
	Percent        uint8
}

type line struct {
	Number uint32 `xml:"number,attr"`
	Hit    uint8  `xml:"hits,attr"`
}

func extractXML(file io.Reader, c chan<- Package) error {
	defer close(c)

	var data xmlFile

	err := xml.NewDecoder(file).Decode(&data)
	if err != nil {
		return err
	}

	for _, d := range data.Packages.Packages {
		p := getPackage(d)

		c <- p
	}

	return err
}

func getPackage(d xmlPackage) Package {
	p := Package{
		Path: strings.ReplaceAll(d.Name, ".", "/"),
	}

	for _, class := range d.Classes.Classes {
		percent, uncovered := class.lines()

		p.Files = append(p.Files, Coverage{
			File:           class.Name,
			Percent:        percent,
			UncoveredLines: uncovered,
		})
	}

	return p
}

func ExtractXMLCommand(_ context.Context, path string, c chan<- Package) error {
	file := filepath.Join(path, "coverage.xml")

	content, err := os.Open(file)
	if err != nil {
		close(c)
		return err
	}

	return extractXML(content, c)
}
