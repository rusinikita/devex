package dashboard

import (
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func heatmap(barNames []string, data values) components.Charter {
	Revert(barNames)

	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Code changes per month",
			Subtitle: "This code changes frequently and a lot",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: false}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			Name: "Date",
			Data: data.timeValues(),
			AxisLabel: &opts.AxisLabel{
				Show:         true,
				ShowMinLabel: true,
				ShowMaxLabel: true,
			},
			SplitArea: &opts.SplitArea{Show: true},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type: "category",
			Name: "Package",
			Data: barNames,
			AxisLabel: &opts.AxisLabel{
				Show:         true,
				Interval:     "0",
				ShowMinLabel: true,
				ShowMaxLabel: true,
			},
			SplitArea: &opts.SplitArea{Show: true},
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Min:        0,
			Max:        float32(data.max()),
			InRange: &opts.VisualMapInRange{
				Color: []string{"#a7d8de", "#eac736", "#d94e5d"},
			},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: strconv.Itoa(200+len(barNames)*30) + "px",
		}),
		charts.WithGridOpts(opts.Grid{
			ContainLabel: true,
		}),
	)

	heatmapdata := Map(data.bar3dValues(), func(in [3]any) opts.HeatMapData {
		return opts.HeatMapData{Value: in}
	})

	hm.AddSeries("Code lines", heatmapdata,
		charts.WithLabelOpts(opts.Label{
			Show:     true,
			Position: "inside",
		}),
	)

	return hm
}
