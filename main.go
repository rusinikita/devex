package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"devex_dashboard/dashboard"
	"devex_dashboard/datacollector"
	"devex_dashboard/datasource"
	"devex_dashboard/datasource/files"
	"devex_dashboard/db"
	"devex_dashboard/project"
)

var tags = flag.String("tags", "", "file content tags")
var lang = flag.String("lang", "go", "main project language")

func main() {
	flag.Parse()

	command := flag.Arg(0)
	alias := flag.Arg(1)

	data := db.DB()

	if command == "" {
		log.Println("no command => running server")
		command = "new"
	}

	switch command {
	case "new":
		path := flag.Arg(2)

		p := project.Project{
			Alias:      alias,
			Language:   *lang,
			FolderPath: path,
			CreatedAt:  time.Now(),
		}

		log.Println("Creating project in", path)

		projectResult := data.FirstOrCreate(&p, p)
		err := projectResult.Error
		if err != nil {
			log.Fatal("db error", err)
		}

		if projectResult.RowsAffected == 0 {
			log.Fatal(alias, " already exists, please use 'update' command")
		}

		if len(*tags) > 0 {
			files.Tags = append(files.Tags, strings.Split(*tags, ",")...)
		}

		err = datacollector.Collect(context.TODO(), data, p, datasource.NewExtractors())
		if err != nil {
			log.Fatal("collect error", err)
		}

	case "update":
		log.Fatal("not implemented, please remove devex_db")

		// p := project.Project{Alias: alias}

		// err := data.Take(&p).Error
		// if err != nil {
		// 	log.Fatal("db error", err)
		// }
		//
		// err = data.Delete(p, p).Error
		// if err != nil {
		// 	log.Fatal("db error", err)
		// }
		//
		// p.ID = 0
		// err = data.Create(&p).Error
		// if err != nil {
		// 	log.Fatal("db error", err)
		// }

		// err = datacollector.Collect(context.TODO(), data, p, datasource.NewExtractors())
		// if err != nil {
		// 	log.Fatal("collect error", err)
		// }

	case "server":
		err := dashboard.RunServer(data)
		if err != nil {
			log.Fatal("server", err)
		}
	}
}
