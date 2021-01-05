package main

import (
	"github.com/hjertnes/roam/commands/create"
	"github.com/hjertnes/roam/commands/diagnostic"
	"github.com/hjertnes/roam/commands/edit"
	"github.com/hjertnes/roam/commands/find"
	"github.com/hjertnes/roam/commands/help"
	iinit "github.com/hjertnes/roam/commands/init"
	"github.com/hjertnes/roam/commands/migrate"
	"github.com/hjertnes/roam/commands/publish"
	"github.com/hjertnes/roam/commands/report"
	"github.com/hjertnes/roam/commands/stats"
	"github.com/hjertnes/roam/commands/sync"
	"github.com/hjertnes/roam/utils"
	"os"
)

func getCreateCommand(path string) *create.Create{
	c, err := create.New(path)
	utils.ErrorHandler(err)

	return c
}

func main() {
	path := utils.GetPath()

	if len(os.Args) == 1 {
		help.Run()

		os.Exit(0)
	}

	switch os.Args[1] {
	case "init":
		utils.ErrorHandler(iinit.Run(path))
	case "publish":
		utils.ErrorHandler(publish.Run(path))
	case "diagnostic":
		utils.ErrorHandler(diagnostic.Run(path))
	case "edit":
		utils.ErrorHandler(edit.Run(path))
	case "migrate":
		utils.ErrorHandler(migrate.Run(path))
	case "sync":
		utils.ErrorHandler(sync.Run(path))
	case "find":
		utils.ErrorHandler(find.Run(path))
	case "create":
		utils.ErrorHandler(getCreateCommand(path).CreateFile(os.Args[2]))
	case "import":
		utils.ErrorHandler(getCreateCommand(path).Import(os.Args[2]))
	case "report":
		utils.ErrorHandler(report.Run(path))
	case "daily":
		utils.ErrorHandler(getCreateCommand(path).Run())
	case "stats":
		utils.ErrorHandler(stats.Run(path))
	default:
		help.Run()
	}
}
