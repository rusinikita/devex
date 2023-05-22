package dashboard

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func treeMap(data values) components.Charter {
	tm := charts.NewTreeMap()

	tm.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeMacarons}),
		charts.WithTitleOpts(opts.Title{
			Title:    "File size chart",
			Subtitle: "Project code lines count",
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:   true,
			Orient: "horizontal",
			Left:   "right",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show: true, Title: "Save as image"},
				Restore: &opts.ToolBoxFeatureRestore{
					Show: true, Title: "Reset"},
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Formatter: opts.FuncOpts(toolTipFormatter),
		}),
	)

	mapOpts := charts.WithTreeMapOpts(
		opts.TreeMapChart{
			Animation:  true,
			Roam:       true,
			UpperLabel: &opts.UpperLabel{Show: true},
			Levels: &[]opts.TreeMapLevel{
				{ // Series
					ItemStyle: &opts.ItemStyle{
						BorderColor: "#777",
						BorderWidth: 1,
						GapWidth:    1,
					},
					UpperLabel: &opts.UpperLabel{Show: false},
				},
				{ // Level
					ItemStyle: &opts.ItemStyle{
						BorderColor: "#666",
						BorderWidth: 2,
						GapWidth:    1,
					},
					Emphasis: &opts.Emphasis{
						ItemStyle: &opts.ItemStyle{BorderColor: "#555"},
					},
				},
				{ // Node
					ColorSaturation: []float32{0.35, 0.5},
					ItemStyle: &opts.ItemStyle{
						GapWidth:              1,
						BorderWidth:           0,
						BorderColorSaturation: 0.6,
					},
				},
			},
		},
	)

	maps := data.simpleMap()
	tm.AddSeries("code", maps,
		mapOpts,
		charts.WithItemStyleOpts(opts.ItemStyle{BorderColor: "#fff"}),
		charts.WithLabelOpts(opts.Label{Show: true, Position: "inside", Color: "White"}),
	)

	return tm
}

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
