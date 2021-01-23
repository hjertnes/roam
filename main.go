package main

import (
	"github.com/hjertnes/roam/commandfactory"
	"github.com/hjertnes/roam/commands/synclog"
	"github.com/hjertnes/roam/utils"
	utilslib "github.com/hjertnes/utils"

	"io/ioutil"
	"os"
	"strings"

	"github.com/hjertnes/roam/commands/configedit"
	"github.com/hjertnes/roam/commands/help"
	iinit "github.com/hjertnes/roam/commands/init"
	"github.com/hjertnes/roam/commands/migrate"
	"github.com/hjertnes/roam/commands/publish"
	"github.com/hjertnes/roam/commands/stats"
	"github.com/hjertnes/roam/commands/version"
)



func getPath() string {
	if utilslib.FileExist("./.roam") {
		data, err := ioutil.ReadFile("./.roam")
		if err == nil {
			return utilslib.ExpandTilde(strings.ReplaceAll(string(data), "\n", ""))
		}
	}

	path, isSet := os.LookupEnv("ROAM")
	if !isSet{
		return utilslib.ExpandTilde("~/txt/roam")
	}

	return utilslib.ExpandTilde(path)
}

func main() {
	path := getPath()

	if len(os.Args) == 1 {
		help.Run(os.Args)

		return
	}

	switch os.Args[1] {
	case "clear":
		utils.ErrorHandler(commandfactory.ClearCommand(path).Run())
	case "init":
		utils.ErrorHandler(iinit.Run(path))
	case "publish":
		utils.ErrorHandler(publish.Run(path, os.Args))
	case "diagnostic":
		utils.ErrorHandler(commandfactory.DiagnosticCommand(path).Run())
	case "config":
		utils.ErrorHandler(configedit.Run(path))
	case "migrate":
		utils.ErrorHandler(migrate.Run(path, os.Args))
	case "sync":
		utils.ErrorHandler(commandfactory.SyncCommand(path).Run())
	case "find":
		utils.ErrorHandler(commandfactory.FindCommand(path).Run())
	case "create":
		utils.ErrorHandler(commandfactory.CreateCommand(path).CreateFile())
	case "import":
		utils.ErrorHandler(commandfactory.CreateCommand(path).RunImport())
	case "daily":
		utils.ErrorHandler(commandfactory.CreateCommand(path).RunDaily())
	case "report":
		utils.ErrorHandler(commandfactory.ReportCommand(path).Run())

	case "stats":
		utils.ErrorHandler(stats.Run(path, os.Args))
	case "log":
		utils.ErrorHandler(synclog.Run(path, os.Args))
	case "version":
		version.Run()
	case "template":
		utils.ErrorHandler(commandfactory.TemplateCommand(path).Run())
	default:
		help.Run(os.Args)
	}
}
