package dashboard

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/rusinikita/devex/internal/helper"
)

func TreeMap(data Values) components.Charter {
	tm := charts.NewTreeMap()

	tm.SetGlobalOptions(
		helper.GetTitleOpts("File size chart", "Project code lines count"),
		helper.GetToolboxOpts(&opts.ToolBoxFeatureRestore{
			Show: true, Title: "Reset"}),
		helper.GetTooltipOptsWithFormatter(),
	)

	maps := data.simpleMap()
	tm.AddSeries("code", maps,
		helper.GetMapOpts(),
		charts.WithItemStyleOpts(opts.ItemStyle{BorderColor: "#fff"}),
		charts.WithLabelOpts(opts.Label{Show: true, Position: "inside", Color: "White"}),
	)

	return tm
}
