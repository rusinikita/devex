package dashboard

import (
	"bytes"
	_ "embed"
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
	ProjectIDs []project.ID `form:"project_ids"`
	PerFiles   bool         `form:"per_files"`
	Filter     string       `form:"filter"`
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

	barNames, data, err := gitChangesData(db, params.PerFiles, dataProjects, params.Filter)
	if err != nil {
		return err
	}

	// RENDER
	page := components.NewPage()

	page.AddCharts(heatmap(barNames, data))
	page.AddCharts(bar3D(barNames, data))

	sizes, err := fileSizes(db, dataProjects, params.Filter)
	if err != nil {
		return err
	}

	page.AddCharts(treeMap(sizes))

	fixes, err := commitMessages(db, params.PerFiles, dataProjects, params.Filter)
	if err != nil {
		return err
	}

	page.AddCharts(bar(fixes))

	contibs, err := contribution(db, dataProjects, params.Filter)
	if err != nil {
		return err
	}

	page.AddCharts(sandkey(contibs))

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
		Filter           string
	}{
		Projects:         projects,
		SelectedProjects: ToSet(params.ProjectIDs),
		PerFiles:         params.PerFiles,
		Filter:           params.Filter,
	}

	err = template.Must(template.New("new").Parse(form)).Execute(&tpl, formData)
	if err != nil {
		return err
	}

	templates.PageTpl = strings.ReplaceAll(templates.PageTpl, "<body>", "<body class=\"container\">\n"+tpl.String())

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
