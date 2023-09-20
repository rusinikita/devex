package middleware

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type UseDatabase struct {
	App      *gin.Engine
	Database *gorm.DB
}

func (m UseDatabase) Use() {
	m.App.Use(func(ctx *gin.Context) {
		ctx.Set("db", m.Database)
	})
}
