package main

import (
	"fmt"
	"github.com/hjertnes/roam/commands/synclog"
	"github.com/hjertnes/roam/errs"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
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
	if !isSet {
		return utilslib.ExpandTilde("~/txt/roam")
	}

	return utilslib.ExpandTilde(path)
}

func errorHandler(err error) {
	if err != nil {
		if eris.Is(err, errs.ErrNotFound) {
			fmt.Println("No matches to search query")
		}

		fmt.Println("Error")

		fmt.Println(eris.ToString(err, true))

		os.Exit(0)
	}
}

func main() {
	path := getPath()

	if len(os.Args) == 1 {
		help.Run(os.Args)

		return
	}

	switch os.Args[1] {
	case "clear":
		errorHandler(buildClearCommand(path).Run())
	case "init":
		errorHandler(iinit.Run(path))
	case "publish":
		errorHandler(publish.Run(path, os.Args))
	case "diagnostic":
		errorHandler(buildDiagnosticCommand(path).Run())
	case "config":
		errorHandler(configedit.Run(path))
	case "migrate":
		errorHandler(migrate.Run(path, os.Args))
	case "sync":
		errorHandler(buildSyncCommand(path).Run())
	case "find":
		errorHandler(buildFindCommand(path).Run())
	case "create":
		errorHandler(buildCreateCommand(path).CreateFile())
	case "import":
		errorHandler(buildCreateCommand(path).RunImport())
	case "report":
		errorHandler(buildReportCommand(path).Run())
	case "daily":
		errorHandler(buildCreateCommand(path).Run())
	case "stats":
		errorHandler(stats.Run(path, os.Args))
	case "log":
		errorHandler(synclog.Run(path, os.Args))
	case "version":
		version.Run()
	case "template":
		errorHandler(buildTemplateCommand(path).Run())
	default:
		help.Run(os.Args)
	}
}
