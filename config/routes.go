package config

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/rusinikita/devex/controller/pages"
	"github.com/rusinikita/devex/middleware"
)

func SetUpMiddlewares(app *gin.Engine, db *gorm.DB) {
	var middlewares []middleware.Middleware
	middlewares = append(middlewares, middleware.UseDatabase{App: app, Database: db})
	middlewares = append(middlewares, middleware.NormalizeError{App: app})

	for _, m := range middlewares {
		m.Use()
	}
}

func SetupRoutes(app *gin.Engine) {
	app.GET("/", pages.Index)
}
