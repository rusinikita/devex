package testcoverage

import (
	"context"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
)

type Package struct {
	Path  string
	Files []Coverage
}

type Coverage struct {
	File           string
	Percent        uint8
	UncoveredLines []uint32
}

type xmlFile struct {
	XMLName  xml.Name `xml:"coverage"`
	Packages struct {
		Packages []xmlPackage `xml:"package"`
	} `xml:"packages"`
}

type xmlPackage struct {
	Name    string `xml:"name,attr"`
	Classes struct {
		Classes []xmlClass `xml:"class"`
	} `xml:"classes"`
}

type xmlClass struct {
	Name  string `xml:"name,attr"`
	Lines struct {
		Lines []line `xml:"line"`
	} `xml:"lines"`
}

func (c xmlClass) lines() (percent uint8, uncovered []uint32) {
	var covered uint32 = 0

	for _, l := range c.Lines.Lines {
		if l.Hit > 0 {
			covered++
		} else {
			uncovered = append(uncovered, l.Number)
		}
	}

	lines := covered + uint32(len(uncovered))
	if lines > 0 {
		percent = uint8((100 * covered) / lines)
	}

	return percent, uncovered
}

type line struct {
	Number uint32 `xml:"number,attr"`
	Hit    uint8  `xml:"hits,attr"`
}

func extractXml(file io.Reader, c chan<- Package) error {
	defer close(c)

	var data xmlFile

	err := xml.NewDecoder(file).Decode(&data)
	if err != nil {
		return err
	}

	for _, d := range data.Packages.Packages {
		p := Package{
			Path: d.Name,
		}

		for _, class := range d.Classes.Classes {
			percent, uncovered := class.lines()

			p.Files = append(p.Files, Coverage{
				File:           class.Name,
				Percent:        percent,
				UncoveredLines: uncovered,
			})
		}

		c <- p
	}

	return err
}

func ExtractXml(_ context.Context, path string, c chan<- Package) error {
	file := filepath.Join(path, "coverage.xml")

	content, err := os.Open(file)
	if err != nil {
		return err
	}

	return extractXml(content, c)
}
