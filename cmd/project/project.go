package project

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"

	"github.com/rusinikita/devex/dao"
	"github.com/rusinikita/devex/datacollector"
	"github.com/rusinikita/devex/datasource"
	"github.com/rusinikita/devex/datasource/files"
	"github.com/rusinikita/devex/db"
)

type Command struct {
	Database *gorm.DB
}

func New(alias string, path string, lang *string, tags *string) datacollector.ErrStruct {
	database := db.DB()

	p, errorStruct := saveProject(database, alias, lang, path)

	if errorStruct.Error != nil {
		return errorStruct
	}

	if len(*tags) > 0 {
		files.Tags = append(files.Tags, strings.Split(*tags, ",")...)
	}

	err := datacollector.Collect(context.TODO(), database, p, datasource.NewExtractors())
	if err != nil {
		return datacollector.ErrStruct{
			Template: alias + "collect error %s",
			Error:    err,
		}
	}
	return datacollector.ErrStruct{}
}

func saveProject(database *gorm.DB, alias string, lang *string, path string) (dao.Project, datacollector.ErrStruct) {
	log.Println("Creating project in", path)
	p := dao.Project{
		Alias:      alias,
		Language:   *lang,
		FolderPath: path,
		CreatedAt:  time.Now(),
	}

	projectResult := database.FirstOrCreate(&p, p)
	err := projectResult.Error
	if err != nil {
		return p, datacollector.ErrStruct{Template: "db error %s", Error: err}
	}

	if projectResult.RowsAffected == 0 {
		err = errors.New("already exists, please use 'update' command")
		return p, datacollector.ErrStruct{Template: alias + " %s", Error: err}
	}

	return p, datacollector.ErrStruct{}
}

func CheckStyle(alias string, path string) datacollector.ErrStruct {
	database := db.DB()

	log.Printf("start parsing \n")

	err := datacollector.CheckStyle(database, alias, path)

	return datacollector.ErrStruct{
		Template: "parsing error %s \n",
		Error:    err,
	}
}

func NotImplemented() datacollector.ErrStruct {
	return datacollector.ErrStruct{
		Template: "%s",
		Error:    errors.New("panic command not set; This is impossible"),
	}
}

func Version() {
	println("v0.1")
}
