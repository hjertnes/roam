package commandfactory

import (
	"github.com/hjertnes/roam/commands/bulkimport"
	"github.com/hjertnes/roam/commands/clear"
	"github.com/hjertnes/roam/commands/createfile"
	"github.com/hjertnes/roam/commands/daily"
	"github.com/hjertnes/roam/commands/diagnostic"
	"github.com/hjertnes/roam/commands/find"
	"github.com/hjertnes/roam/commands/publish"
	"github.com/hjertnes/roam/commands/report"
	"github.com/hjertnes/roam/commands/sync"
	"github.com/hjertnes/roam/commands/template"
	"github.com/hjertnes/roam/utils"
	"os"
)

func CreateFileCommand(path string) *createfile.CreateFile {
	c, err := createfile.New(path, os.Args)
	utils.ErrorHandler(err)

	return c
}

func ImportCommand(path string) *bulkimport.Import {
	c, err := bulkimport.New(path, os.Args)
	utils.ErrorHandler(err)

	return c
}

func PublishCommand(path string) *publish.Publish {
	c, err := publish.New(path, os.Args)
	utils.ErrorHandler(err)

	return c
}

func DailyCommand(path string) *daily.Daily {
	c, err := daily.New(path, os.Args)
	utils.ErrorHandler(err)

	return c
}

func ClearCommand(path string) *clear.Clear {
	c, err := clear.New(path, os.Args)
	utils.ErrorHandler(err)

	return c
}

func FindCommand(path string) *find.Find {
	c, err := find.New(path, os.Args)
	utils.ErrorHandler(err)

	return c
}

func TemplateCommand(path string) *template.Template{
	t, err := template.New(path, os.Args)
	utils.ErrorHandler(err)

	return t
}

func SyncCommand(path string) *sync.Sync{
	t, err := sync.New(path, os.Args)
	utils.ErrorHandler(err)

	return t
}

func ReportCommand(path string) *report.Report{
	t, err := report.New(path, os.Args)
	utils.ErrorHandler(err)

	return t
}

func DiagnosticCommand(path string) *diagnostic.Diagnostic{
	t, err := diagnostic.New(path, os.Args)
	utils.ErrorHandler(err)

	return t
}
