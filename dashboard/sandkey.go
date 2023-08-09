package dashboard

import (
	"path"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/rusinikita/devex/slices"
)

func sandkey(data values) components.Charter {
	sk := charts.NewSankey()

	var (
		nodeNames []string
		links     []opts.SankeyLink
	)

	projects := map[string]float32{}

	for _, v := range data {
		p := path.Join(v.Alias, v.Package, v.Name)
		nodeNames = append(nodeNames, v.Author, v.Alias, p)

		projects[path.Join(v.Alias, v.Author)] += float32(v.Value)

		links = append(links, opts.SankeyLink{
			Source: v.Author,
			Target: p,
			Value:  float32(v.Value),
		})
	}

	for p, value := range projects {
		links = append(links, opts.SankeyLink{
			Source: path.Dir(p),
			Target: path.Base(p),
			Value:  value,
		})
	}

	nodes := slices.Map(slices.Distinct(nodeNames), func(in string) opts.SankeyNode {
		return opts.SankeyNode{
			Name: in,
			ItemStyle: &opts.ItemStyle{
				Color: colorful.FastWarmColor().Hex(),
			},
		}
	})

	sk.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: strconv.Itoa(200+len(nodes)*20) + "px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Contribution in last year",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "item",
			TriggerOn: "mousemove|click",
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:   true,
			Orient: "horizontal",
			Left:   "right",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show: true, Title: "Save as image"},
			},
		}),
	)

	sk.AddSeries("sankey", nodes, links,
		charts.WithLineStyleOpts(opts.LineStyle{
			Color:     "source",
			Curveness: 0.5,
		}),
		charts.WithLabelOpts(opts.Label{Show: true}),
	)

	return sk
}
