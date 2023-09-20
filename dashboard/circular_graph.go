package dashboard

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/rusinikita/devex/internal/helper"
)

const (
	minLines = 200
	maxLines = 500
)

const (
	minHeight = 500
	maxHeight = 1000
)

const countNodes = 40
const curvenessCircularGraph = 0.3

func CircularGraph(data AllImports) components.Charter {
	sk := charts.NewGraph()

	categories, nodes, links := data.tree()
	minLinesLabel, height := getSizesByNodes(nodes)

	sk.SetGlobalOptions(
		helper.GetInitializationOpts(height),
		helper.GetTitleOpts("Dependencies", "Lines - connections between modules, node size - module code lines"),
		helper.GetTooltipOptsCalcGraph(),
		helper.GetToolboxOpts(&opts.ToolBoxFeatureRestore{Show: true, Title: "Reset"}),
	)

	sk.AddSeries("graph", nodes, links,
		helper.GetChartOpts(categories),
		helper.GetLineStyleOpts("target", curvenessCircularGraph),
		helper.GetLabelOptsWithFormatter(minLinesLabel),
		helper.GetItemStyleOpts(),
	)

	return sk
}

func getSizesByNodes(nodes []opts.GraphNode) (minLinesLabel int, height int) {
	minLinesLabel = minLines
	height = minHeight
	if len(nodes) > countNodes {
		minLinesLabel = maxLines
		height = maxHeight
	}

	return minLinesLabel, height
}
