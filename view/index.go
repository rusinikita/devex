package view

import (
	"bytes"
	_ "embed"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/templates"

	"github.com/rusinikita/devex/dao"
	"github.com/rusinikita/devex/dashboard"
	"github.com/rusinikita/devex/project"
	"github.com/rusinikita/devex/slices"
)

//go:embed form.gohtml
var form string

const countFiles = 20

type requestDtoStruct struct {
	Prefs        []string
	Data         dashboard.Values
	HeatmapBars  dashboard.Values
	Projects     []dao.Project
	DataProjects []project.ID
	SQLFilter    string
	FilesTop     dashboard.Values
}

func RenderPage(db *gorm.DB, params Params, w http.ResponseWriter) error {
	// request
	requestDto, err := transformRequestToDto(db, params)
	if err != nil {
		return err
	}

	// get page
	page, err := getPage(db, params, requestDto)
	if err != nil {
		return err
	}

	// response
	return page.Render(w)
}

func getPage(db *gorm.DB, params Params, dto requestDtoStruct) (*components.Page, error) { //nolint
	page := components.NewPage()
	packagePrefs := dto.Prefs
	dataProjects := dto.DataProjects
	sqlFilter := dto.SQLFilter
	projects := dto.Projects

	// start block first
	// Code changes per month; Top changes speed; file size chart
	err := addTopChangesToPage(db, page, dto)
	if err != nil {
		return page, err
	}
	// end block first

	// start block second
	fileTagsData, err := addCommitsToPage(db, page, params, dto)
	if err != nil {
		return page, err
	}
	// end block second

	// start block third
	// TODO I can’t get into the function because... there will be 5 arguments

	barDto := dashboard.BarDto{
		Name:   "File tags",
		Desc:   fmt.Sprintf("Files with tags from '%s' filter in file content", params.FileFilters),
		Values: fileTagsData,
	}

	page.AddCharts(dashboard.Bar(barDto))

	contibs, err := dashboard.Contribution(db, params.PerFiles, dataProjects, sqlFilter)
	if err != nil {
		return page, err
	}

	page.AddCharts(dashboard.Sandkey(contibs.WithPackagesTrimmed(packagePrefs)))

	fileImports, err := dashboard.Imports(db, params.PerFilesImports, dataProjects, sqlFilter)
	if err != nil {
		return page, err
	}

	page.AddCharts(dashboard.CircularGraph(fileImports.WithPackagesTrimmed(packagePrefs)))

	page.SetLayout(components.PageNoneLayout)
	page.AddCustomizedCSSAssets("https://cdn.jsdelivr.net/npm/@picocss/pico@1/css/pico.min.css")
	// end block third

	// template hack
	originTpl := templates.PageTpl
	defer func() { templates.PageTpl = originTpl }()

	var tpl bytes.Buffer
	formData := struct {
		SelectedProjects slices.Set[project.ID]
		PackageFilter    string
		NameFilter       string
		TrimPackage      string
		CommitFilters    string
		FileFilters      string
		Projects         []dao.Project
		PerFiles         bool
		PerFilesImports  bool
	}{
		Projects:         projects,
		SelectedProjects: slices.ToSet(params.ProjectIDs),
		PerFiles:         params.PerFiles,
		PerFilesImports:  params.PerFilesImports,
		PackageFilter:    params.PackageFilter,
		NameFilter:       params.NameFilter,
		TrimPackage:      params.TrimPackage,
		CommitFilters:    params.CommitFilters,
		FileFilters:      params.FileFilters,
	}

	err = template.Must(template.New("new").Parse(form)).Execute(&tpl, formData)
	if err != nil {
		return page, err
	}

	templates.PageTpl = strings.ReplaceAll(templates.PageTpl, "<body>", "<body class=\"container\">\n"+tpl.String())
	templates.PageTpl = strings.ReplaceAll(templates.PageTpl, "<html>", "<html data-theme=\"light\">")

	return page, nil
}

func addCommitsToPage( //nolint
	db *gorm.DB,
	page *components.Page,
	params Params,
	dto requestDtoStruct,
) (dashboard.Values, error) {
	sqlFilter := dto.SQLFilter
	dataProjects := dto.DataProjects
	packagePrefs := dto.Prefs

	commitFilter := sqlFilter + " and " + slices.SQLFilter("c.message", params.CommitFilters)
	fileCommits, err := dashboard.CommitMessages(db, params.PerFiles, dataProjects, commitFilter)
	if err != nil {
		return nil, err
	}

	fileCommits = fileCommits.WithPackagesTrimmed(packagePrefs)

	barDto := dashboard.BarDto{
		Name:   "Commits",
		Desc:   fmt.Sprintf("Changes with '%s' filter applied to file", params.CommitFilters),
		Values: fileCommits,
	}
	page.AddCharts(dashboard.Bar(barDto))

	tagsFilter := " and " + slices.SQLFilter("tags", params.FileFilters)
	fileTagsData, err := dashboard.FileTags(db, dataProjects, sqlFilter, tagsFilter)
	if err != nil {
		return nil, err
	}

	fileTagsData = fileTagsData.TagsToValue(params.FileFilters).WithPackagesTrimmed(packagePrefs)

	return fileTagsData, nil
}

func addTopChangesToPage(db *gorm.DB, page *components.Page, dto requestDtoStruct) error {
	charts, err := getCharts(db, dto)

	if err != nil {
		return err
	}

	for i := 0; i < len(charts); i++ {
		page.AddCharts(charts[i])
	}

	return nil
}

func getCharts(db *gorm.DB, dto requestDtoStruct) (charts []components.Charter, err error) {
	packagePrefs := dto.Prefs

	barNames := dto.HeatmapBars.WithPackagesTrimmed(packagePrefs).BarNames()
	barNamesData := dto.Data.WithPackagesTrimmed(packagePrefs)
	barDto := getBarDto(dto.FilesTop)

	sizes, err := dashboard.FileSizes(db, dto.DataProjects, dto.SQLFilter)
	if err != nil {
		return charts, err
	}

	trimmedSizes := sizes.WithPackagesTrimmed(packagePrefs)

	return []components.Charter{
		dashboard.Heatmap(barNames, barNamesData),
		dashboard.Bar(barDto),
		dashboard.TreeMap(trimmedSizes),
	}, nil
}

func getBarDto(top dashboard.Values) dashboard.BarDto {
	return dashboard.BarDto{
		Name:   "Top changes speed",
		Desc:   "List of packages/files ordered by average change lines per month speed",
		Values: top,
	}
}

func transformRequestToDto(db *gorm.DB, params Params) (requestDtoStruct, error) { //nolint
	var projects []dao.Project

	if err := db.Find(&projects).Error; err != nil {
		return requestDtoStruct{}, err
	}

	dataProjects := getDataProjects(params.ProjectIDs, projects)

	sqlFilter := params.sqlFilter()
	packagePrefs := strings.Split(params.TrimPackage, ",")

	filesTop, err := dashboard.GitChangesTop(db, params.PerFiles, dataProjects, sqlFilter)
	if err != nil {
		return requestDtoStruct{}, err
	}

	heatmapBars := getHeatMapBars(filesTop)

	data, err := dashboard.GitChangesData(db, params.PerFiles, dataProjects, heatmapBars)
	if err != nil {
		return requestDtoStruct{}, err
	}

	return requestDtoStruct{
		Prefs:        packagePrefs,
		Data:         data,
		HeatmapBars:  heatmapBars,
		Projects:     projects,
		DataProjects: dataProjects,
		SQLFilter:    sqlFilter,
		FilesTop:     filesTop,
	}, nil
}

func getHeatMapBars(filesTop dashboard.Values) dashboard.Values {
	if len(filesTop) > countFiles {
		return filesTop[:countFiles]
	}

	return filesTop
}

func getDataProjects(projectIDs []project.ID, projects []dao.Project) []project.ID {
	dataProjects := projectIDs

	if len(projectIDs) == 0 {
		for _, p := range projects {
			dataProjects = append(dataProjects, p.ID)
		}
	}

	return dataProjects
}
