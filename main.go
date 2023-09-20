package main

import (
	"flag"
	"log"
	"os"

	"github.com/rusinikita/devex/cmd/dashboard"
	"github.com/rusinikita/devex/cmd/project"
	"github.com/rusinikita/devex/datacollector"
	"github.com/rusinikita/devex/internal/constants"
)

var tags = flag.String("tags", "", "file content tags")
var lang = flag.String("lang", "go", "main project language")

func main() {
	commandName, alias := getFlags()

	errStruct := execute(commandName, alias)

	if errStruct.Error != nil {
		log.Printf(errStruct.Template, errStruct.Error)
		os.Exit(constants.ERROR)
	}

	log.Printf("parsing success \n")
	os.Exit(constants.SUCCESS)
}

func execute(name string, alias string) datacollector.ErrStruct {
	errStruct := datacollector.ErrStruct{}

	switch name {
	case "new":
		path := flag.Arg(constants.THIRD)
		errStruct = project.New(alias, path, lang, tags)
	case "server":
		errStruct = dashboard.RunServer()
	case "version":
		project.Version()
	case "check_style":
		path := flag.Arg(constants.THIRD)
		errStruct = project.CheckStyle(alias, path)
	case "update":
	default:
		errStruct = project.NotImplemented()
	}

	return errStruct
}

func getFlags() (commandName string, alias string) {
	flag.Parse()

	commandName = getCommand()
	alias = flag.Arg(constants.SECOND)

	return commandName, alias
}

func getCommand() string {
	command := flag.Arg(constants.FIRST)
	if command == "" {
		log.Println("no command => running server")
		command = "new"
	}
	return command
}
