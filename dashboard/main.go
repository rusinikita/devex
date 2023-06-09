package dashboard

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/templates"
	"gorm.io/gorm"

	"devex_dashboard/project"
)

//go:embed form.gohtml
var form string

type Params struct {
	ProjectIDs      []project.ID `form:"project_ids"`
	PerFiles        bool         `form:"per_files"`
	PerFilesImports bool         `form:"per_files_imports"`
	PackageFilter   string       `form:"package_filter"`
	NameFilter      string       `form:"name_filter"`
	TrimPackage     string       `form:"trim_package"`
	CommitFilters   string       `form:"commit_filters"`
	FileFilters     string       `form:"file_filters"`
}

func (p Params) sqlFilter() (sql string) {
	if p.PackageFilter != "" {
		sql += "and " + SQLFilter("package", p.PackageFilter)
	}

	if p.NameFilter != "" {
		sql += " and " + SQLFilter("name", p.NameFilter)
	}

	return sql
}

func renderPage(db *gorm.DB, params Params, w http.ResponseWriter) error {
	// REQUEST
	var projects []project.Project
	if err := db.Find(&projects).Error; err != nil {
		return err
	}

	dataProjects := params.ProjectIDs
	if len(params.ProjectIDs) == 0 {
		for _, p := range projects {
			dataProjects = append(dataProjects, p.ID)
		}
	}

	sqlFilter := params.sqlFilter()
	packagePrefs := strings.Split(params.TrimPackage, ",")

	barNames, data, err := gitChangesData(db, params.PerFiles, dataProjects, sqlFilter)
	if err != nil {
		return err
	}

	// RENDER
	page := components.NewPage()

	page.AddCharts(heatmap(barNames.withPackagesTrimmed(packagePrefs).barNames(), data.withPackagesTrimmed(packagePrefs)))

	sizes, err := fileSizes(db, dataProjects, sqlFilter)
	if err != nil {
		return err
	}

	page.AddCharts(treeMap(sizes.withPackagesTrimmed(packagePrefs)))

	fileCommits, err := commitMessages(db, params.PerFiles, dataProjects, sqlFilter, " and "+SQLFilter("c.message", params.CommitFilters))
	if err != nil {
		return err
	}

	fileCommits = fileCommits.withPackagesTrimmed(packagePrefs)

	page.AddCharts(bar("Commits", fmt.Sprintf("Changes with '%s' filter applied to file", params.CommitFilters), fileCommits))

	fileTagsData, err := fileTags(db, dataProjects, sqlFilter, " and "+SQLFilter("tags", params.FileFilters))
	if err != nil {
		return err
	}

	fileTagsData = fileTagsData.tagsToValue(params.FileFilters).withPackagesTrimmed(packagePrefs)

	page.AddCharts(bar("File tags", fmt.Sprintf("Files with tags from '%s' filter in file content", params.FileFilters), fileTagsData))

	// fileContents, err := commitMessages(db, params.PerFiles, dataProjects, sqlFilter, " and "+SQLFilter("c.message", params.CommitFilters))
	// if err != nil {
	// 	return err
	// }
	//
	// page.AddCharts(bar("Contents", "Files with keywords in content",fileContents))

	contibs, err := contribution(db, params.PerFiles, dataProjects, sqlFilter)
	if err != nil {
		return err
	}

	page.AddCharts(sandkey(contibs.withPackagesTrimmed(packagePrefs)))

	fileImports, err := imports(db, params.PerFilesImports, dataProjects, sqlFilter)
	if err != nil {
		return err
	}

	page.AddCharts(circularGraph(fileImports.withPackagesTrimmed(packagePrefs)))

	page.SetLayout(components.PageNoneLayout)
	page.AddCustomizedCSSAssets("https://cdn.jsdelivr.net/npm/@picocss/pico@1/css/pico.min.css")

	// template hack
	originTpl := templates.PageTpl
	defer func() { templates.PageTpl = originTpl }()

	var tpl bytes.Buffer
	formData := struct {
		Projects         []project.Project
		SelectedProjects Set[project.ID]
		PerFiles         bool
		PerFilesImports  bool
		PackageFilter    string
		NameFilter       string
		TrimPackage      string
		CommitFilters    string
		FileFilters      string
	}{
		Projects:         projects,
		SelectedProjects: ToSet(params.ProjectIDs),
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
		return err
	}

	templates.PageTpl = strings.ReplaceAll(templates.PageTpl, "<body>", "<body class=\"container\">\n"+tpl.String())
	templates.PageTpl = strings.ReplaceAll(templates.PageTpl, "<html>", "<html data-theme=\"light\">")

	return page.Render(w)
}

func RunServer(db *gorm.DB) error {
	engine := gin.New()

	engine.Use(func(ctx *gin.Context) {
		ctx.Next()

		err := ctx.Errors.Last()
		if err == nil {
			return
		}

		code := http.StatusBadRequest
		if err.IsType(gin.ErrorTypePrivate) {
			code = http.StatusInternalServerError
		}

		ctx.JSON(code, err.JSON())
	})

	engine.GET("/", func(ctx *gin.Context) {
		params := Params{}

		if err := ctx.BindQuery(&params); err != nil {
			return
		}

		err := renderPage(db, params, ctx.Writer)
		if err != nil {
			ctx.Error(err)
		}
	})

	return engine.Run(":1080")
}
