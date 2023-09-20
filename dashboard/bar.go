package dashboard

import (
	"path"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/rusinikita/devex/internal/helper"
	"github.com/rusinikita/devex/slices"
)

type BarDto struct {
	Name   string
	Desc   string
	Values Values
}

func Bar(dto BarDto) components.Charter {
	dtoName := dto.Name
	desc := dto.Desc
	data := dto.Values
	slices.Revert(data)

	names := data.BarNames()

	barData := getBarData(data)

	bar := getBar(dtoName, desc, names, barData)

	bar.AddSeries("", barData)

	return bar
}

func getBarData(data Values) []opts.BarData {
	var barData []opts.BarData
	for _, v := range data {
		barData = append(barData, opts.BarData{
			Name:  path.Join(v.Alias, v.Package, v.Name),
			Value: v.Value,
		})
	}

	return barData
}

func getBar(dtoName string, desc string, names []string, barData []opts.BarData) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		helper.GetTitleOpts(dtoName, desc),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		helper.GetYAxisOpts(names),
		helper.GetXAxisOpts(),
		helper.GetGridOpts(),
		helper.GetInitializationOpts(helper.DefaultHeight+len(barData)*20),
		helper.GetToolboxOpts(&opts.ToolBoxFeatureRestore{}),
	)

	return bar
}
