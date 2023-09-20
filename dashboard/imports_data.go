package dashboard

import (
	"path"
	"strings"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/rusinikita/devex/slices"
)

type importsData struct {
	Alias   string
	Package string
	Name    string
	Imports []string `gorm:"serializer:json"`
	Lines   uint32
}

type AllImports []importsData

type allImportsStruct struct {
	MaxLines int
}

func (all AllImports) WithPackagesTrimmed(prefixes []string) AllImports {
	for i := range all {
		all[i].Package = slices.MultiTrimPrefix(all[i].Package, prefixes)

		for k := range all[i].Imports {
			all[i].Imports[k] = slices.MultiTrimPrefix(all[i].Imports[k], prefixes)
		}
	}

	return all
}

func (all AllImports) tree() ([]*opts.GraphCategory, []opts.GraphNode, []opts.GraphLink) {
	projectTrees, importsStruct := all.getTrees()

	return importsStruct.processingTrees(projectTrees)
}

func (all AllImports) getTrees() (map[string]*file, allImportsStruct) {
	projectTrees := map[string]*file{}

	maxLines := 0

	for _, data := range all {
		// Is there already a key in the hash map? If not, add it
		project := getProjectByAliasOrSet(projectTrees, data.Alias)

		// is there a key for child paths in the hash map, if not add
		projectFile := getFileByDataOrSet(project, data)

		// determine the largest number of lines in files, assigning the number of lines to a file
		maxLines = setValueAndMaxLinesByFile(projectFile, data.Lines, maxLines)

		setChildrens(projectFile, data.Imports)
	}

	return projectTrees, allImportsStruct{
		MaxLines: maxLines,
	}
}

func setChildrens(projectFile *file, imports []string) {
	for _, imprt := range imports {
		if imprt == "" || len(imprt) < countImports {
			continue
		}
		// the logic apparently comes from some kind of language
		if !strings.Contains(imprt, "/") {
			imprt = strings.ReplaceAll(imprt, ".", "/")
		}

		// TODO Why add it if we are removing it afterwards?
		projectFile.children[imprt] = nil
	}
}

func setValueAndMaxLinesByFile(projectFile *file, dataLines uint32, lines int) int {
	moduleLines := projectFile.value + int(dataLines)
	projectFile.value = moduleLines
	if moduleLines > lines {
		return moduleLines
	}

	return lines
}

func getFileByDataOrSet(project *file, data importsData) *file {
	filePath := path.Join(data.Package, strings.TrimSuffix(data.Name, ".py"))
	f, ok := project.children[filePath]

	if !ok {
		f = newFile(filePath, 0)
		project.children[filePath] = f
	}

	return f
}

func getProjectByAliasOrSet(projectTrees map[string]*file, alias string) *file {
	project, ok := projectTrees[alias]

	if !ok {
		project = newFile(alias, 0)
		projectTrees[alias] = project
	}

	return project
}

func (importsStruct allImportsStruct) processingTrees(
	trees map[string]*file,
) ([]*opts.GraphCategory, []opts.GraphNode, []opts.GraphLink) {
	var (
		categories []*opts.GraphCategory
		nodes      []opts.GraphNode
		links      []opts.GraphLink
	)

	for alias, project := range trees {
		// Some kind of private crutch, apparently for some kind of logic
		if alias == "push" {
			alias = "_push_"
		}

		categories = appendCategory(alias, categories)

		nodes, links = importsStruct.processingChildren(project.children, alias, nodes, links)
	}

	return categories, nodes, links
}

func (importsStruct allImportsStruct) processingChildren(
	children map[string]*file,
	alias string,
	nodes []opts.GraphNode,
	links []opts.GraphLink,
) ([]opts.GraphNode, []opts.GraphLink) {
	resultNodes := nodes
	resultLinks := links

	for fPath, f := range children {
		node := getGraphNode(importsStruct, float32(f.value), alias, fPath)

		resultNodes = append(resultNodes, node)

		resultLinks = getResultLinks(resultLinks, f.children, node.Name, alias)
	}

	return resultNodes, resultLinks
}

func getResultLinks(links []opts.GraphLink, children map[string]*file, name string, alias string) []opts.GraphLink {
	resultLinks := links

	for fileImport := range children {
		// TODO Why add it if we are removing it afterwards?
		_, hasImport := children[fileImport]
		if !hasImport {
			continue
		}

		resultLinks = append(resultLinks, getLink(name, alias, fileImport))
	}

	return resultLinks
}

func getLink(name string, alias string, fileImport string) opts.GraphLink {
	return opts.GraphLink{
		Source: name,
		Target: path.Join(alias, fileImport),
	}
}

func getGraphNode(importsStruct allImportsStruct, fValue float32, alias string, fPath string) opts.GraphNode {
	return opts.GraphNode{
		Name:       path.Join(alias, fPath),
		SymbolSize: fValue * percent100 / float32(importsStruct.MaxLines),
		Value:      fValue,
		Category:   alias,
		ItemStyle: &opts.ItemStyle{
			Color: colorful.FastWarmColor().Hex(),
		},
	}
}

func appendCategory(alias string, categories []*opts.GraphCategory) []*opts.GraphCategory {
	return append(categories, &opts.GraphCategory{
		Name: alias,
		Label: &opts.Label{
			Show: true,
		},
	})
}
