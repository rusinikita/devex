/*
Package search is WIP experiment with semantic code search.

Transfer code files to vectors using transformers, saving vectors to database.
Then allow to search for code places by meaning/prompt or code sample.

Current problems:
- It can't be integrated in files data collection pipeline because it is too slow (5 seconds per 4000 symbols file);
- SQLite vss (vector search) extension requires complex build pipeline and I won't use it.
Alternative is using external vector DB (such as milvus, qdrant).
- Search request prototype implemented though table full scan
*/
package search

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/acheong08/vectordb/rank"
	"github.com/acheong08/vectordb/vectors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"gorm.io/gorm"

	database "github.com/rusinikita/devex/db"
	"github.com/rusinikita/devex/project"
	"github.com/rusinikita/devex/slices"
)

type query struct {
	Query   string `form:"q"`
	Results []result
}

type result struct {
	project.File
	project.Project
}

//go:embed form.gohtml
var form string

func Hander(ctx *gin.Context) {
	db := database.GetDB(ctx)

	q := query{}

	if err := ctx.BindQuery(&q); err != nil {
		return
	}

	if q.Query == "" {
		ctx.Render(http.StatusOK, render.HTML{
			Template: template.Must(template.New("search").Parse(form)),
			Name:     "search",
			Data:     q,
		})

		return
	}

	qVec, err := vectors.Encode(q.Query)
	if err != nil {
		ctx.Error(err)

		return
	}

	var topFiles []project.File

	var files []project.File
	err = db.Where("present > 0").FindInBatches(&files, 1000, func(tx *gorm.DB, batch int) error {
		topFiles = append(topFiles, topK(qVec, files, 5)...)

		return nil
	}).Error
	if err != nil {
		ctx.Error(err)

		return
	}

	topFiles = topK(qVec, topFiles, 10)

	ids := slices.Map(files, func(f project.File) project.ID {
		return f.ID
	})

	err = db.Model(result{}).Table("files").Where("id in ?", ids).
		Joins("join projects p on p.id = files.project").
		Scan(&q.Results).Error
	if err != nil {
		ctx.Error(err)

		return
	}

	ctx.Render(http.StatusOK, render.HTML{
		Template: template.Must(template.New("search").Parse(form)),
		Name:     "search",
		Data:     q,
	})
}

func topK(query []float64, files []project.File, k int) (topFiles []project.File) {
	corpus := slices.Map(files, func(f project.File) []float64 {
		return nil // f.SemanticVectors[0]
	})

	result := rank.Rank([][]float64{query}, corpus, k, false)[0]

	for _, searchResult := range result {
		topFiles = append(topFiles, files[searchResult.CorpusID])
	}

	return topFiles
}
