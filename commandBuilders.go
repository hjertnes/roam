package main

import (
	"github.com/hjertnes/roam/commands/clear"
	"github.com/hjertnes/roam/commands/create"
	"github.com/hjertnes/roam/commands/diagnostic"
	"github.com/hjertnes/roam/commands/find"
	"github.com/hjertnes/roam/commands/report"
	"github.com/hjertnes/roam/commands/sync"
	"github.com/hjertnes/roam/commands/template"
	"os"
)

func buildCreateCommand(path string) *create.Create {
	c, err := create.New(path, os.Args)
	errorHandler(err)

	return c
}

func buildClearCommand(path string) *clear.Clear {
	c, err := clear.New(path, os.Args)
	errorHandler(err)

	return c
}

func buildFindCommand(path string) *find.Find {
	c, err := find.New(path, os.Args)
	errorHandler(err)

	return c
}

func buildTemplateCommand(path string) *template.Template{
	t, err := template.New(path, os.Args)
	errorHandler(err)

	return t
}

func buildSyncCommand(path string) *sync.Sync{
	t, err := sync.New(path, os.Args)
	errorHandler(err)

	return t
}

func buildReportCommand(path string) *report.Report{
	t, err := report.New(path, os.Args)
	errorHandler(err)

	return t
}

func buildDiagnosticCommand(path string) *diagnostic.Diagnostic{
	t, err := diagnostic.New(path, os.Args)
	errorHandler(err)

	return t
}
