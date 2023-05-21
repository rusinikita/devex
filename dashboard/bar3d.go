package dashboard

import (
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func bar3D(barNames []string, data timeSeriesData) components.Charter {
	var bar3dNames []string
	copy(bar3dNames, barNames)
	sort.Strings(bar3dNames)

	bar3d := charts.NewBar3D()
	bar3d.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Code changes per month",
			Subtitle: "This code changes frequently and a lot",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: false}),
		charts.WithXAxis3DOpts(opts.XAxis3D{
			Type: "category",
			Name: "Date",
			Data: data.timeValues(),
		}),
		charts.WithYAxis3DOpts(opts.YAxis3D{
			Type: "category",
			Name: "Package",
			Data: bar3dNames,
		}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{
			Name: "Lines",
			Type: "value",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width: "100%",
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Min:        0,
			Max:        float32(data.max()),
			InRange: &opts.VisualMapInRange{
				Color: []string{"#a7d8de", "#eac736", "#d94e5d"},
			},
		}),
	)

	chart3d := Map(data.bar3dValues(), func(in [3]any) opts.Chart3DData {
		return opts.Chart3DData{Value: in[:]}
	})

	bar3d.AddSeries("Code lines", chart3d,
		charts.WithBar3DChartOpts(opts.Bar3DChart{
			Shading: "lambert",
		}),
	)

	return bar3d
}
