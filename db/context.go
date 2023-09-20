package db

import (
	"context"

	"gorm.io/gorm"
)

func GetDB(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
