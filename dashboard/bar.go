package dashboard

import (
	"fmt"
	"path"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/rusinikita/devex/slices"
)

func bar(name, desc string, data values) components.Charter {
	slices.Revert(data)

	names := data.barNames()

	var barData []opts.BarData
	for _, v := range data {
		barData = append(barData, opts.BarData{
			Name:  path.Join(v.Alias, v.Package, v.Name),
			Value: v.Value,
		})
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    name,
			Subtitle: desc,
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Package",
			Type: "category",
			Show: true,
			Data: names,
			AxisLabel: &opts.AxisLabel{
				Show:         true,
				ShowMinLabel: true,
				ShowMaxLabel: true,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Show: true,
			Name: "Fixes count",
			Type: "value",
		}),
		charts.WithGridOpts(opts.Grid{
			ContainLabel: true,
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: fmt.Sprintf("%dpx", 200+20*len(barData)),
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

	bar.AddSeries("", barData)

	return bar
}
