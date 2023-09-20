package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NormalizeError struct {
	App *gin.Engine
}

func (m NormalizeError) Use() {
	m.App.Use(func(ctx *gin.Context) {
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
}
