package commands

import (
	"context"
	"fmt"
	"github.com/hjertnes/roam/models"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

func printListOfLinks(output []string, links []models.File) []string{
	if len(links) == 0 {
		output = append(output, "No links")
	}

	for _, link := range links {
		output = append(output, fmt.Sprintf("- %s\n", link.Path))
	}

	return output
}

func Report(path string) error {
	conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
	if err != nil {
		return eris.Wrap(err, "failed to get config")
	}

	ctx := context.Background()

	pxp, err := pgxpool.Connect(ctx, conf.DatabaseConnectionString)
	if err != nil {
		return eris.Wrap(err, "could not connect to database")
	}

	dal := dal2.New(ctx, pxp)

	files, err := dal.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	output := make([]string, 0)

	for _, file := range files {
		output = append(output, fmt.Sprintf("# %s", file.Path))

		links, err := dal.GetLinks(file.ID)
		if err != nil {
			return eris.Wrap(err, "failed to get list of links")
		}

		backlinks, err := dal.GetBacklinks(file.ID)
		if err != nil {
			return eris.Wrap(err, "failed to get list of backlinks")
		}

		output = append(output, "## Links")

		output = printListOfLinks(output, links)

		output = append(output, "## Backlinks")

		output = printListOfLinks(output, backlinks)
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	out, err := r.Render(strings.Join(output, "\n"))
	if err != nil {
		return eris.Wrap(err, "failed to render markdown file")
	}

	fmt.Print(out)

	return nil
}
