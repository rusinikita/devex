package helper

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const DefaultHeight = 200

const toolTipFormatter = `
function (info) {
	var formatUtil = echarts.format;
	var value = info.value;
	var treePathInfo = info.treePathInfo;
	var treePath = [];
	for (var i = 1; i < treePathInfo.length; i++) {
		treePath.push(treePathInfo[i].name);
	}
	return ['<div class="tooltip-title">' + formatUtil.encodeHTML(treePath.join('/')) + '</div>',
		'Size: ' + formatUtil.addCommas(value) + ' lines',
		].join('');
}
`

func GetTitleOpts(title string, subtitle string) charts.GlobalOpts {
	return charts.WithTitleOpts(opts.Title{
		Title:    title,
		Subtitle: subtitle,
	})
}

func GetTooltipOptsWithFormatter() charts.GlobalOpts {
	return charts.WithTooltipOpts(opts.Tooltip{
		Show:      true,
		Formatter: opts.FuncOpts(toolTipFormatter),
	})
}

func GetTooltipOpts() charts.GlobalOpts {
	return charts.WithTooltipOpts(opts.Tooltip{
		Show: true,
	})
}

func GetToolboxOpts(restore *opts.ToolBoxFeatureRestore) charts.GlobalOpts {
	return charts.WithToolboxOpts(opts.Toolbox{
		Show:   true,
		Orient: "horizontal",
		Left:   "right",
		Feature: &opts.ToolBoxFeature{
			SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
				Show: true, Title: "Save as image"},
			Restore: restore,
		},
	})
}

func GetInitializationOpts(height int) charts.GlobalOpts {
	return charts.WithInitializationOpts(opts.Initialization{
		Width:  "100%",
		Height: fmt.Sprintf("%dpx", height),
	})
}

func GetGridOpts() charts.GlobalOpts {
	return charts.WithGridOpts(opts.Grid{
		ContainLabel: true,
	})
}

func GetXAxisOpts() charts.GlobalOpts {
	return charts.WithXAxisOpts(opts.XAxis{
		Show: true,
		Name: "Fixes count",
		Type: "value",
	})
}

func GetYAxisOpts(names []string) charts.GlobalOpts {
	return charts.WithYAxisOpts(opts.YAxis{
		Name: "Package",
		Type: "category",
		Show: true,
		Data: names,
		AxisLabel: &opts.AxisLabel{
			Show:         true,
			ShowMinLabel: true,
			ShowMaxLabel: true,
		},
	})
}

func GetVisualMapOpts(max float32) charts.GlobalOpts {
	return charts.WithVisualMapOpts(opts.VisualMap{
		Calculable: true,
		Min:        0,
		Max:        max,
		InRange: &opts.VisualMapInRange{
			Color: []string{"#a7d8de", "#eac736", "#d94e5d"},
		},
	})
}

func GetYAxisOptsHeatMap(names []string) charts.GlobalOpts {
	return charts.WithYAxisOpts(opts.YAxis{
		Type: "category",
		Name: "Package",
		Data: names,
		AxisLabel: &opts.AxisLabel{
			Show:         true,
			Interval:     "0",
			ShowMinLabel: true,
			ShowMaxLabel: true,
		},
		SplitArea: &opts.SplitArea{Show: true},
	})
}

func GetXAxisOptsHeatMap(data any) charts.GlobalOpts {
	return charts.WithXAxisOpts(opts.XAxis{
		Type: "category",
		Name: "Date",
		Data: data,
		AxisLabel: &opts.AxisLabel{
			Show:         true,
			ShowMinLabel: true,
			ShowMaxLabel: true,
		},
		SplitArea: &opts.SplitArea{Show: true},
	})
}

func GetLegendOpts() charts.GlobalOpts {
	return charts.WithLegendOpts(opts.Legend{Show: false})
}

func GetTooltipOptsSandKey() charts.GlobalOpts {
	return charts.WithTooltipOpts(opts.Tooltip{
		Show:      true,
		Trigger:   "item",
		TriggerOn: "mousemove|click",
	})
}

func GetTooltipOptsCalcGraph() charts.GlobalOpts {
	return charts.WithTooltipOpts(opts.Tooltip{
		Show:      true,
		Trigger:   "item",
		TriggerOn: "mousemove|click",
		Formatter: "{b}: {c} lines",
	})
}
