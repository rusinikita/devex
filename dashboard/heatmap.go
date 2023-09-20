package dashboard

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/rusinikita/devex/internal/helper"
	"github.com/rusinikita/devex/slices"
)

const countInnerHeatMapData = 3

func Heatmap(barNames []string, data Values) components.Charter {
	slices.Revert(barNames)

	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		helper.GetTitleOpts("Code changes per month", "This code changes frequently and a lot"),
		helper.GetTooltipOpts(),
		helper.GetLegendOpts(),
		helper.GetToolboxOpts(&opts.ToolBoxFeatureRestore{}),
		helper.GetXAxisOptsHeatMap(data.TimeValues()),
		helper.GetYAxisOptsHeatMap(barNames),
		helper.GetVisualMapOpts(float32(data.max())),
		helper.GetInitializationOpts(helper.DefaultHeight+len(barNames)*30),
		helper.GetGridOpts(),
	)

	heatmapData := getHeatMapData(data.bar3dValues())

	hm.AddSeries("Code lines", heatmapData, helper.GetLabelOpts("inside"))

	return hm
}

func getHeatMapData(values [][countInnerHeatMapData]any) []opts.HeatMapData {
	return slices.Map(values, func(in [countInnerHeatMapData]any) opts.HeatMapData {
		return opts.HeatMapData{Value: in}
	})
}
