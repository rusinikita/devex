package dashboard

import (
	"path"

	"github.com/go-echarts/go-echarts/v2/opts"
)

type file struct {
	children map[string]*file
	name     string
	value    int
}

func newFile(name string, value int) *file {
	return &file{
		name:     name,
		value:    value,
		children: map[string]*file{},
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
