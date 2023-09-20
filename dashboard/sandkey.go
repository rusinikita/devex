package dashboard

import (
	"path"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/rusinikita/devex/internal/helper"
	"github.com/rusinikita/devex/slices"
)

const curvenessSandKey = 0.5

func Sandkey(data Values) components.Charter {
	sk := charts.NewSankey()

	nodeNames, links := getLinks(data)
	nodes := getNodes(nodeNames)

	sk.SetGlobalOptions(
		helper.GetInitializationOpts(helper.DefaultHeight+len(nodes)*20),
		helper.GetTitleOpts("Contribution in last year", ""),
		helper.GetTooltipOptsSandKey(),
		helper.GetToolboxOpts(&opts.ToolBoxFeatureRestore{}),
	)

	sk.AddSeries("sankey", nodes, links,
		helper.GetLineStyleOpts("source", curvenessSandKey),
		helper.GetLabelOpts(""),
	)

	return sk
}

func getLinks(data Values) (nodeNames []string, links []opts.SankeyLink) {
	projects := map[string]float32{}

	for _, v := range data {
		p := path.Join(v.Alias, v.Package, v.Name)
		nodeNames = append(nodeNames, v.Author, v.Alias, p)
		floatValue := float32(v.Value)

		projects[path.Join(v.Alias, v.Author)] += floatValue

		links = appendLink(links, v.Author, p, floatValue)
	}

	for p, value := range projects {
		links = appendLink(links, path.Dir(p), path.Base(p), value)
	}

	return nodeNames, links
}

func appendLink(links []opts.SankeyLink, author string, target string, value float32) []opts.SankeyLink {
	return append(links, opts.SankeyLink{
		Source: author,
		Target: target,
		Value:  value,
	})
}

func getNodes(names []string) []opts.SankeyNode {
	return slices.Map(slices.Distinct(names), func(in string) opts.SankeyNode {
		return opts.SankeyNode{
			Name: in,
			ItemStyle: &opts.ItemStyle{
				Color: colorful.FastWarmColor().Hex(),
			},
		}
	})
}
