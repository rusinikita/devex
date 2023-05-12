package main

import (
	"context"
	"flag"
	"log"
	"time"

	"devex_dashboard/dashboard"
	"devex_dashboard/datacollector"
	"devex_dashboard/datasource"
	"devex_dashboard/db"
	"devex_dashboard/project"
)

func main() {
	flag.Parse()

	command := flag.Arg(0)
	alias := flag.Arg(1)

	data := db.DB()

	switch command {
	case "new":
		path := flag.Arg(2)
		p := project.Project{
			Alias:      alias,
			Language:   "python",
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

		err = datacollector.Collect(context.TODO(), data, p, datasource.NewExtractors())
		if err != nil {
			log.Fatal("collect error", err)
		}

	case "update":
		p := project.Project{Alias: alias}

		err := data.Take(&p).Error
		if err != nil {
			log.Fatal("db error", err)
		}

		log.Fatal("update not implemented yet")

		// err = datacollector.Collect(context.TODO(), data, p, datasource.NewExtractors())
		// if err != nil {
		// 	log.Fatal("collect error", err)
		// }

	case "server":
		err := dashboard.RunServer(context.TODO(), data)
		if err != nil {
			log.Fatal("server", err)
		}
	}
}
