package helper

import "github.com/go-echarts/go-echarts/v2/opts"

const (
	saturation35 = 0.35
	saturation50 = 0.5
	saturation60 = 0.6
)

const (
	width1 = 2
	width2 = 2
)

func GetNode() opts.TreeMapLevel {
	return opts.TreeMapLevel{
		ColorSaturation: []float32{saturation35, saturation50},
		ItemStyle: &opts.ItemStyle{
			GapWidth:              1,
			BorderWidth:           0,
			BorderColorSaturation: saturation60,
		},
	}
}

func GetLevel() opts.TreeMapLevel {
	return opts.TreeMapLevel{
		ItemStyle: &opts.ItemStyle{
			BorderColor: "#666",
			BorderWidth: width2,
			GapWidth:    width1,
		},
		Emphasis: &opts.Emphasis{
			ItemStyle: &opts.ItemStyle{BorderColor: "#555"},
		},
	}
}

func GetSeries() opts.TreeMapLevel {
	return opts.TreeMapLevel{
		ItemStyle: &opts.ItemStyle{
			BorderColor: "#777",
			BorderWidth: width1,
			GapWidth:    width1,
		},
		UpperLabel: &opts.UpperLabel{Show: false},
	}
}
