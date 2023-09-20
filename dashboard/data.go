package dashboard

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/rusinikita/devex/slices"
)

type ValueData struct {
	Tags    map[string]uint32 `gorm:"serializer:json"`
	Alias   string
	Package string
	Name    string
	Author  string
	Time    string
	Value   float64
}

type Values []ValueData

func (v Values) Len() int           { return len(v) }
func (v Values) Less(i, j int) bool { return v[i].Value > v[j].Value }
func (v Values) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

const countInnerSlice = 3
const countResults = 20
const countImports = 3
const percent100 = 100

func (v Values) WithPackagesTrimmed(prefixes []string) Values {
	for i := range v {
		v[i].Package = slices.MultiTrimPrefix(v[i].Package, prefixes)
	}

	return v
}

func (v Values) BarNames() []string {
	return slices.Map(v, func(d ValueData) string {
		return filepath.Join(d.Alias, d.Package, d.Name)
	})
}

func (v Values) TimeValues() (r []string) {
	for _, d := range v {
		r = append(r, d.Time)
	}

	r = slices.Distinct(r)
	sort.Strings(r)

	return r
}

func (v Values) bar3dValues() (r [][countInnerSlice]any) {
	return slices.Map(v, func(d ValueData) [countInnerSlice]any {
		nameFormat := filepath.Join(d.Alias, d.Package, d.Name)

		return [3]any{d.Time, nameFormat, d.Value}
	})
}

func (v Values) max() float64 {
	return slices.Fold(v, func(item ValueData, value float64) float64 {
		if item.Value > value {
			return item.Value
		}

		return value
	})
}

func (v Values) simpleMap() (result []opts.TreeMapNode) {
	root := newFile("", 0)

	for _, data := range v {
		setChildrensToTM(data, root)
	}

	return root.treeNode().Children
}

func setChildrensToTM(data ValueData, root *file) {
	project, ok := root.children[data.Alias]
	if !ok {
		project = newFile(data.Alias, 0)
		root.children[data.Alias] = project
	}

	folder, ok := project.children[data.Package]
	if !ok {
		folder = newFile(data.Package, 0)
		project.children[data.Package] = folder
	}

	dashboardFile := newFile(data.Name, int(data.Value))

	folder.children[data.Name] = dashboardFile
}

func (v Values) TagsToValue(tagsFilter string) (result Values) {
	tags := map[string]bool{}

	for _, split := range strings.Split(tagsFilter, ";") {
		for _, s := range strings.Split(split, ",") {
			tags[strings.TrimSuffix(s, "!")] = true
		}
	}

	return getValues(v, tags)
}

func getValues(v Values, tags map[string]bool) Values {
	result := Values{}

	for _, value := range v {
		data := value

		for tag, count := range data.Tags {
			if tags[tag] {
				data.Value += float64(count)
			}
		}

		result = append(result, data)
	}

	return GetResult(result)
}

func GetResult(result Values) Values {
	sort.Sort(result)

	if len(result) > countResults {
		return result[:countResults]
	}

	return result
}
