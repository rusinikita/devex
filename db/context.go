package db

import (
	"context"

	"gorm.io/gorm"
)

func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, "db", db)
}

func GetDB(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
