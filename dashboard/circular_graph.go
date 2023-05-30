package dashboard

import (
	"fmt"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func circularGraph(data allImports) components.Charter {
	sk := charts.NewGraph()

	categories, nodes, links := data.tree()

	minLinesLabel := 200
	height := 500
	if len(nodes) > 40 {
		height = 1000
		minLinesLabel = 500
	}

	sk.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: strconv.Itoa(height) + "px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Dependencies",
			Subtitle: "Lines - connections between modules, node size - module code lines",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "item",
			TriggerOn: "mousemove|click",
			Formatter: "{b}: {c} lines",
		}),
	)

	sk.AddSeries("graph", nodes, links,
		charts.WithGraphChartOpts(opts.GraphChart{
			Layout:             "circular",
			Categories:         categories,
			Roam:               true,
			FocusNodeAdjacency: true,
			Draggable:          true,
		}),
		charts.WithLineStyleOpts(opts.LineStyle{
			Color:     "target",
			Curveness: 0.3,
		}),
		charts.WithLabelOpts(opts.Label{
			Show: true,
			Formatter: opts.FuncOpts(fmt.Sprintf(`function (info) {
	return info.value > %d ? info.name : '';
}`, minLinesLabel)),
		}),
		charts.WithItemStyleOpts(opts.ItemStyle{
			GapWidth: 100,
		}),
	)

	return sk
}
