package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"devex_dashboard/datasource"
)

func createDB(file string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(datasource.DataEntities()...)
	if err != nil {
		panic(err)
	}

	// deleteTX := db.Session(&gorm.Session{AllowGlobalUpdate: true})
	// for _, e := range datasource.DataEntities() {
	// 	deleteTX.Delete(e)
	// }

	return db
}

func DB() *gorm.DB {
	return createDB("devex_bd.db")
}

func TestDB() *gorm.DB {
	return createDB("file::memory:?cache=shared")
}
