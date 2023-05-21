package dashboard

import (
	"path/filepath"
	"sort"
)

type valueData struct {
	Alias   string
	Package string
	Name    string
	Value   float64
}

type values []valueData

func (v values) barNames() []string {
	return Map(v, func(d valueData) string {
		return filepath.Join(d.Alias, d.Package, d.Name)
	})
}

type timedData struct {
	Alias   string
	Package string
	Name    string
	BarTime string
	Value   int64
}

type timeSeriesData []timedData

func (l timeSeriesData) timeValues() (r []string) {
	for _, d := range l {
		r = append(r, d.BarTime)
	}

	r = Distinct(r)
	sort.Strings(r)

	return r
}

func (l timeSeriesData) bar3dValues() (r [][3]any) {
	return Map(l, func(d timedData) [3]any {
		nameFormat := filepath.Join(d.Alias, d.Package, d.Name)

		return [3]any{d.BarTime, nameFormat, d.Value}
	})
}

func (l timeSeriesData) max() int64 {
	return Fold(l, func(item timedData, value int64) int64 {
		if item.Value > value {
			return item.Value
		}

		return value
	})
}
