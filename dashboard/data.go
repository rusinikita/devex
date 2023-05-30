package dashboard

import (
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/lucasb-eyer/go-colorful"
)

type valueData struct {
	Alias   string
	Package string
	Name    string
	Author  string
	Value   float64
}

type values []valueData

func (v values) barNames() []string {
	return Map(v, func(d valueData) string {
		return filepath.Join(d.Alias, d.Package, d.Name)
	})
}

func (v values) treeMaps() (result []opts.TreeMapNode) {
	root := newFile("", 0)

	for _, data := range v {
		p := append([]string{data.Alias}, strings.Split(data.Package, "/")...)

		root.insert(p, data.Name, int(data.Value))
	}

	return root.treeNode().Children
}

func (v values) simpleMap() (result []opts.TreeMapNode) {
	root := newFile("", 0)

	for _, data := range v {
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

		folder.children[data.Name] = newFile(data.Name, int(data.Value))
	}

	return root.treeNode().Children
}

type file struct {
	name     string
	children map[string]*file
	value    int
}

func newFile(name string, value int) *file {
	f := &file{
		name:  name,
		value: value,
	}

	if value == 0 {
		f.children = map[string]*file{}
	}

	return f
}

func (f *file) insert(filePath []string, name string, value int) {
	for _, folderName := range filePath {
		folder, ok := f.children[folderName]
		if !ok {
			folder = newFile(folderName, 0)
			f.children[folderName] = folder
		}

		f = folder
	}

	f.children[name] = &file{
		name:  name,
		value: value,
	}
}

func (f file) treeNode() opts.TreeMapNode {
	node := opts.TreeMapNode{
		Name:     f.name,
		Value:    f.value,
		Children: nil,
	}

	for _, ff := range f.children {
		if len(f.children) == 1 {
			folder := node.Name
			node = ff.treeNode()
			node.Name = path.Join(folder, node.Name)
			break
		}

		node.Children = append(node.Children, ff.treeNode())
	}

	return node
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

type importsData struct {
	Alias   string
	Package string
	Name    string
	Lines   uint32
	Imports []string `gorm:"serializer:json"`
}

type allImports []importsData

func (all allImports) tree() (categories []*opts.GraphCategory, nodes []opts.GraphNode, links []opts.GraphLink) {
	projectTrees := map[string]*file{}

	maxLines := 0

	for _, data := range all {
		project, ok := projectTrees[data.Alias]
		if !ok {
			project = newFile(data.Alias, 0)
			projectTrees[data.Alias] = project
		}

		filePath := path.Join(strings.TrimPrefix(data.Package, "src/"), strings.TrimSuffix(data.Name, ".py"))
		f, ok := project.children[filePath]
		if !ok {
			f = newFile(filePath, 0)
			project.children[filePath] = f

		}

		moduleLines := f.value + int(data.Lines)
		f.value = moduleLines
		if moduleLines > maxLines {
			maxLines = moduleLines
		}

		for _, i := range data.Imports {
			trimmedImport := strings.TrimPrefix(i, "go.avito.ru/msg/service-seller-audience/")
			trimmedImport = strings.TrimPrefix(trimmedImport, "go.avito.ru/av/service-messenger-push/")
			if !strings.Contains(trimmedImport, "/") {
				trimmedImport = strings.ReplaceAll(trimmedImport, ".", "/")
			}
			if len(trimmedImport) < 3 || trimmedImport == "internal" {
				continue
			}

			f.children[trimmedImport] = nil
		}
	}

	for alias, project := range projectTrees {
		if alias == "push" {
			alias = "_push_"
		}

		categories = append(categories, &opts.GraphCategory{
			Name: alias,
			Label: &opts.Label{
				Show: true,
			},
		})

		for fPath, f := range project.children {
			node := opts.GraphNode{
				Name:       path.Join(alias, fPath),
				SymbolSize: float32(f.value) * 100 / float32(maxLines),
				Value:      float32(f.value),
				Category:   alias,
				ItemStyle: &opts.ItemStyle{
					Color: colorful.FastWarmColor().Hex(),
				},
			}

			nodes = append(nodes, node)

			for fileImport := range f.children {
				_, hasImport := project.children[fileImport]
				if !hasImport {
					continue
				}

				link := opts.GraphLink{
					Source: node.Name,
					Target: path.Join(alias, fileImport),
				}

				links = append(links, link)
			}
		}
	}

	return categories, nodes, links
}
