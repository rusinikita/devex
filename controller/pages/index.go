package pages

import (
	"github.com/gin-gonic/gin"

	database "github.com/rusinikita/devex/db"
	"github.com/rusinikita/devex/view"
)

func Index(ctx *gin.Context) {
	params := view.Params{}

	if err := ctx.BindQuery(&params); err != nil {
		return
	}

	err := view.RenderPage(database.GetDB(ctx), params, ctx.Writer)
	if err != nil {
		err = ctx.Error(err)
	}

	if err != nil {
		return
	}
}
