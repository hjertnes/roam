package main

import (
	"fmt"
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
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
	"os"
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

func f(path string) {
	err := find.Run(path)
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
	c, err := create.New(path)
	errorHandler(err)

	switch len(os.Args) {
	case two:
		err = c.DailyToday()
		errorHandler(err)
		return
	case three:
		err := c.Daily(os.Args[2])
		errorHandler(err)
	default:
		help.Run()
	}
}

func main() {
	path := utils.GetPath()

	if len(os.Args) == 1 {
		help.Run()

		return
	}

	var err error

	switch os.Args[1] {
	case "init":
		err = iinit.Run(path)
	case "publish":
		to := ""
		if len(os.Args) > 2{
			to = os.Args[2]
			if to == "--include-private"{
				to = ""
			}
		}
		excludePrivate := true

		for _, a := range os.Args{
			if a == "--include-private"{
				excludePrivate = false
			}
		}

		err = publish.Run(path, to, excludePrivate)
	case "diagnostic":
		err = diagnostic.Run(path)
	case "edit":
		err = edit.Run(path)
	case "migrate":
		err = migrate.Run(path)
	case "sync":
		err = sync.Run(path)
	case "find":
		f(path)
	case "create":
		c, err := create.New(path)
		if err != nil{
			break
		}
		err = c.CreateFile(os.Args[2])
	case "import":
		c, err := create.New(path)
		errorHandler(err)
		err = c.Import(os.Args[2])
		errorHandler(err)
	case "report":
		err = report.Run(path)
	case "daily":
		daily(path)
	case "stats":
		err = stats.Run(path)

	default:
		help.Run()
	}

	errorHandler(err)
}
