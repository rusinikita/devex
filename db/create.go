package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/rusinikita/devex/datasource"
)

func createDB(file string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if res := db.Exec("PRAGMA foreign_keys = ON", nil); res.Error != nil {
		panic(res.Error)
	}

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
	return createDB("devex.db")
}

func TestDB(file ...string) *gorm.DB {
	if len(file) > 0 {
		return createDB(file[0])
	}

	return createDB("file::memory:?cache=shared")
}
