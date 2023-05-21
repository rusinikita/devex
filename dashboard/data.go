package dashboard

import (
	"path/filepath"
	"sort"
)

type barData struct {
	Alias   string
	Package string
	Name    string
	BarTime string
	Value   int64
}

type dataList []barData

func (l dataList) timeValues() (r []string) {
	for _, d := range l {
		r = append(r, d.BarTime)
	}

	r = Distinct(r)
	sort.Strings(r)

	return r
}

func (l dataList) bar3dValues() (r [][3]any) {
	return Map(l, func(d barData) [3]any {
		nameFormat := filepath.Join(d.Alias, d.Package, d.Name)

		return [3]any{d.BarTime, nameFormat, d.Value}
	})
}

func (l dataList) max() int64 {
	return Fold(l, func(item barData, value int64) int64 {
		if item.Value > value {
			return item.Value
		}

		return value
	})
}
