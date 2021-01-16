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

	"github.com/hjertnes/roam/commands/clear"
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

func getCreateCommand(path string) *create.Create {
	c, err := create.New(path)
	errorHandler(err)

	return c
}

func getClearCommand(path string) *clear.Clear {
	c, err := clear.New(path)
	errorHandler(err)

	return c
}

func getFindCommand(path string) *find.Find {
	c, err := find.New(path)
	errorHandler(err)

	return c
}

func main() {
	path := getPath()

	if len(os.Args) == 1 {
		help.Run()

		os.Exit(0)
	}

	switch os.Args[1] {
	case "clear":
		errorHandler(getClearCommand(path).Run())
	case "init":
		errorHandler(iinit.Run(path))
	case "publish":
		errorHandler(publish.Run(path))
	case "diagnostic":
		errorHandler(diagnostic.Run(path))
	case "edit":
		errorHandler(edit.Run(path))
	case "migrate":
		errorHandler(migrate.Run(path))
	case "sync":
		errorHandler(sync.Run(path))
	case "find":
		errorHandler(getFindCommand(path).Run())
	case "create":
		errorHandler(getCreateCommand(path).CreateFile())
	case "import":
		errorHandler(getCreateCommand(path).RunImport())
	case "report":
		errorHandler(report.Run(path))
	case "daily":
		errorHandler(getCreateCommand(path).Run())
	case "stats":
		errorHandler(stats.Run(path))
	case "log":
		errorHandler(synclog.Run(path))
	case "version":
		version.Run()
	default:
		help.Run()
	}
}
