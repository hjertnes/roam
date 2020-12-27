package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hjertnes/roam/commands"
	"github.com/hjertnes/roam/configuration"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
)

const (
	two   = 2
	three = 3
)

func errorHandler(err error) {
	if err != nil {
		fmt.Println("Error")
		fmt.Println(eris.ToString(err, true))

		os.Exit(0)
	}
}

func find(path string) {
	err := commands.FindEdit(path)
	if err != nil {
		if eris.Is(err, errs.ErrNotFound) {
			fmt.Println("No matches to search query")

			return
		}

		errorHandler(err)

		return
	}
}

func daily(path string) {
	switch len(os.Args) {
	case two:
		conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
		errorHandler(err)
		err = commands.Daily(path, time.Now().Format(conf.DateFormat))
		errorHandler(err)

		return
	case three:
		err := commands.Daily(path, os.Args[2])
		errorHandler(err)
	default:
		commands.Help()
	}
}

func main() {
	path := utils.GetPath()

	if len(os.Args) == 1 {
		commands.Help()

		return
	}

	var err error

	switch os.Args[1] {
	case "init":
		err = commands.Init(path)
	case "diagnostic":
		err = commands.Diagnostic(path)
	case "edit":
		err = commands.Edit(path)
	case "migrate":
		err = commands.Migrate(path)
	case "sync":
		err = commands.Sync(path)
	case "find":
		find(path)
	case "create":
		err = commands.Create(path, os.Args[2])
	case "report":
		err = commands.Report(path)
	case "daily":
		daily(path)
	case "stats":
		err = commands.Stats(path)

	default:
		commands.Help()
	}

	errorHandler(err)
}
