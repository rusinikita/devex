package testcoverage

import "encoding/xml"

type Packages struct {
	Packages []xmlPackage `xml:"package"`
}

type Classes struct {
	Classes []xmlClass `xml:"class"`
}

type Lines struct {
	Lines []line `xml:"line"`
}

type xmlFile struct {
	XMLName  xml.Name `xml:"coverage"`
	Packages Packages `xml:"packages"`
}

type xmlPackage struct {
	Name    string  `xml:"name,attr"`
	Classes Classes `xml:"classes"`
}

type xmlClass struct {
	Name  string `xml:"name,attr"`
	Lines Lines  `xml:"lines"`
}

func (c xmlClass) lines() (percent uint8, uncovered []uint32) {
	var covered uint32

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
