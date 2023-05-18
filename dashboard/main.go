package dashboard

import (
	"net/http"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/templates"
	"gorm.io/gorm"
)

type barData struct {
	Package string
	Name    string
	BarTime string
	Value   int64
}

func init() {
	templates.PageTpl = strings.ReplaceAll(templates.PageTpl, "</body>\n</html>", "<p>Biba boba</p></body>\n</html>")
}

/*

select package, name, date(`time`, 'start of month') as month, sum(rows_added + rows_removed) as line_changes
from git_changes as ch
         join files f on ch.file = f.id
where f.project = 1
  and f.present > 0
  and package like 'internal%'
  and package not like '%generated%'
  and package not like '%mocks%'
	and package not like '%mocks%' and time > date('now', '-6 month')
group by package, name, date(`time`, 'start of month')
order by date(`time`, 'start of month') desc;
*/

// TOTO firstly get points with lines > median from all points

func data(db *gorm.DB) (barNames []string, result []barData, err error) {
	var bars []struct {
		Package string
		Name    string
	}

	sqlBars := `
	with fcm as (select package, name, date("time", 'start of month') as month, sum(rows_added + rows_removed) as line_changes
	from git_changes as ch
    join files f on ch.file = f.id
		 where f.project = ?
		   and f.present > 0
		   and package like 'internal%'
		   and package not like '%generated%'
		   and package not like '%mocks%'
		   and name not like '%_test.go'
		   and name not like '%mock%'
		   and time > date('now', '-48 month')
		group by package, name, date("time", 'start of month'))
	select package, name, count(*), sum(line_changes), avg(line_changes)
	from fcm group by package, name
	having count(*) > 6
	order by avg(line_changes) desc
	limit 20
`

	err = db.Raw(sqlBars, 1).Scan(&bars).Error
	if err != nil {
		return nil, nil, err
	}

	for _, bar := range bars {
		barNames = append(barNames, path.Join(bar.Package, bar.Name))
	}

	sql := `
	select package, name, date("time", 'start of month') as bar_time, sum(rows_added + rows_removed) as value
	from git_changes as ch
	join files f on ch.file = f.id
	where f.project = ?
		and f.present > 0
		and package || '/' || name in ?
		and time > date('now', '-24 month')
	group by package, name, date("time", 'start of month')
`

	err = db.Raw(sql, 1, barNames).Scan(&result).Error

	return barNames, result, err
}

func httpserver(db *gorm.DB, w http.ResponseWriter) {
	// REQUEST
	barNames, data, err := data(db)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	var (
		timePoints  []string
		maxValue    int64
		heatmapdata []opts.HeatMapData
		chart3d     []opts.Chart3DData
	)

	for _, d := range data {
		timeFormat := d.BarTime
		nameFormat := filepath.Join(d.Package, d.Name)
		timePoints = append(timePoints, timeFormat)
		value := []any{timeFormat, nameFormat, d.Value}
		heatmapdata = append(heatmapdata, opts.HeatMapData{Value: value})
		chart3d = append(chart3d, opts.Chart3DData{Value: value})
		if d.Value > maxValue {
			maxValue = d.Value
		}
	}

	for i := 0; i < len(barNames)/2; i++ {
		barNames[i], barNames[len(barNames)-i-1] = barNames[len(barNames)-i-1], barNames[i]
	}

	timePoints = Distinct(timePoints)
	sort.Strings(timePoints)

	// RENDER
	page := components.NewPage()

	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Code changes per month",
			Subtitle: "This code changes frequently and a lot",
			Left:     "40%",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "10%"}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			Name: "Date",
			Data: timePoints,
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
			Max:        float32(maxValue),
			InRange: &opts.VisualMapInRange{
				Color: []string{"#a7d8de", "#eac736", "#d94e5d"},
			},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "90%",
			Height: strconv.Itoa(len(barNames)*30) + "px",
		}),
		charts.WithGridOpts(opts.Grid{
			ContainLabel: true,
		}),
	)

	hm.AddSeries("Code lines", heatmapdata,
		charts.WithLabelOpts(opts.Label{
			Show:     true,
			Position: "inside",
		}),
	)

	bar3d := charts.NewBar3D()
	bar3d.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Code changes per month",
			Subtitle: "This code changes frequently and a lot",
			Left:     "40%",
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "10%"}),
		charts.WithXAxis3DOpts(opts.XAxis3D{
			Type: "category",
			Name: "Date",
			Data: timePoints,
		}),
		charts.WithYAxis3DOpts(opts.YAxis3D{
			Type: "category",
			Name: "Package",
			Data: barNames,
		}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{
			Name: "Lines",
			Type: "value",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width: "90%",
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Min:        0,
			Max:        float32(maxValue),
			InRange: &opts.VisualMapInRange{
				Color: []string{"#a7d8de", "#eac736", "#d94e5d"},
			},
		}),
		// charts.WithGrid3DOpts(opts.Grid3D{
		// 	BoxWidth: 200,
		// 	BoxDepth: 80,
		// }),
	)

	bar3d.AddSeries("Code lines", chart3d,
		charts.WithBar3DChartOpts(opts.Bar3DChart{
			Shading: "lambert",
		}),
	)

	page.AddCharts(hm)
	page.AddCharts(bar3d)

	page.SetLayout(components.PageCenterLayout)

	errPanic(page.Render(w))
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func RunServer(db *gorm.DB) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		httpserver(db, writer)
	})
	return http.ListenAndServe(":1080", nil)
}

func Distinct[T comparable](a []T) []T {
	hash := make(map[T]struct{})

	for _, v := range a {
		hash[v] = struct{}{}
	}

	set := make([]T, 0, len(hash))
	for k := range hash {
		set = append(set, k)
	}

	// sort.Slice(set, func(i, j int) bool {
	// 	return set[i] < set[j]
	// })

	return set
}
