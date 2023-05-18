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

	// go test -coverprofile=covcov -covermode=atomic ./...
	// file:start_line.start_symbol,end_line.end_symbol numberOfStatements hits
	// считаем hits == 0 и все строки между start_line и end_line
	/*
		mode: atomic
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:30.153,38.2 1 0
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:40.49,53.33 5 0
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:53.33,56.13 2 0
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:56.13,59.8 2 0
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:59.8,61.19 2 0
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:61.19,62.11 1 0
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:63.11,65.6 1 0
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:70.2,72.12 2 0
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:75.125,80.16 4 5
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:80.16,83.3 2 0
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:85.2,89.100 4 5
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:89.100,96.3 4 0
		go.avito.ru/msg/service-seller-audience/internal/actions/send_discounts/action.go:98.2,109.16 2 5
	*/

	err := xml.NewDecoder(file).Decode(&data)
	if err != nil {
		return err
	}

	for _, d := range data.Packages.Packages {
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
