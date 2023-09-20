package helper

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func GetMapOpts() charts.SeriesOpts {
	return charts.WithTreeMapOpts(
		opts.TreeMapChart{
			Animation:  true,
			Roam:       true,
			UpperLabel: &opts.UpperLabel{Show: true},
			Levels: &[]opts.TreeMapLevel{
				GetSeries(),
				GetLevel(),
				GetNode(),
			},
		},
	)
}

func GetLabelOpts(position string) charts.SeriesOpts {
	return charts.WithLabelOpts(opts.Label{
		Show:     true,
		Position: position,
	})
}

func GetLineStyleOpts(color string, curveness float32) charts.SeriesOpts {
	return charts.WithLineStyleOpts(opts.LineStyle{
		Color:     color,
		Curveness: curveness,
	})
}

func GetItemStyleOpts() charts.SeriesOpts {
	return charts.WithItemStyleOpts(opts.ItemStyle{
		GapWidth: 100,
	})
}

func GetLabelOptsWithFormatter(minLinesLabel int) charts.SeriesOpts {
	return charts.WithLabelOpts(opts.Label{
		Show: true,
		Formatter: opts.FuncOpts(fmt.Sprintf(`function (info) {
	return info.value > %d ? info.name : '';
}`, minLinesLabel)),
	})
}

func GetChartOpts(categories []*opts.GraphCategory) charts.SeriesOpts {
	return charts.WithGraphChartOpts(opts.GraphChart{
		Layout:             "circular",
		Categories:         categories,
		Roam:               true,
		FocusNodeAdjacency: true,
		Draggable:          true,
	})
}
