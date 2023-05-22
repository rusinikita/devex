package dashboard

import (
	"path"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func sandkey(data values) components.Charter {
	sk := charts.NewSankey()

	var (
		nodeNames []string
		links     []opts.SankeyLink
	)

	for _, v := range data {
		p := path.Join(v.Alias, v.Package)
		nodeNames = append(nodeNames, v.Name, v.Alias, p)

		links = append(links,
			opts.SankeyLink{
				Source: v.Name,
				Target: p,
				Value:  float32(v.Value),
			},
			opts.SankeyLink{
				Source: p,
				Target: v.Alias,
				Value:  float32(v.Value),
			},
		)
	}

	nodes := Map(Distinct(nodeNames), func(in string) opts.SankeyNode {
		return opts.SankeyNode{Name: in}
	})

	sk.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: strconv.Itoa(len(nodes)*10) + "px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Contribution in last year",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "item",
			TriggerOn: "mousemove|click",
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
