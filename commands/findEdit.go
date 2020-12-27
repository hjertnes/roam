package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/charmbracelet/glamour"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/utils"
	"github.com/hjertnes/roam/widgets/selectinput"
	"github.com/hjertnes/roam/widgets/textinput"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

func FindEdit(path string) error {
	search, err := textinput.Run("Search for a note", "Search: ")
	if err != nil {
		return eris.Wrap(err, "failed to get a search string from textinput")
	}

	fmt.Println("Loading...")

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

	result, err := dal.Find(search)
	if err != nil {
		return eris.Wrap(err, "failed to search for files in database")
	}

	choice, err := selectinput.Run("Select match", filesToChoices(result))
	if err != nil {
		return eris.Wrap(err, "failed to get selection from selectinput")
	}

	action, err := selectinput.Run("Select action", []selectinput.Choice{
		{Title: "Edit", Value: "edit"},
		{Title: "View", Value: "view"},
		{Title: "Backlinks", Value: "backlinks"},
	})
	if err != nil {
		return eris.Wrap(err, "could not get action")
	}

	if action.Value == "edit" {
		editor := utils.GetEditor()
		cmd := exec.Command(editor, choice.Value)

		err = cmd.Start()
		if err != nil {
			return eris.Wrap(err, "could not open file in editor")
		}

		return nil
	}

	if action.Value == "view" {

		data, err := ioutil.ReadFile(choice.Value)
		if err != nil {
			return eris.Wrap(err, "could not read file")
		}
		metadata := models.Fm{}
		err = frontmatter.Unmarshal(data, &metadata)
		if err != nil {
			return eris.Wrap(err, "could not unmarkshal frontmatter")
		}

		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
		)

		out, err := r.Render(fmt.Sprintf("# %s\n%s", metadata.Title, metadata.Content))
		if err != nil {
			return eris.Wrap(err, "failed to render markdown file")
		}
		fmt.Print(out)
		return nil
	}

	if action.Value == "backlinks" {
		var file *models.File
		for _, r := range result {
			if r.Path == choice.Value {
				file = &r
				break
			}
		}

		if file == nil {
			return eris.Wrap(errs.NotFound, "No file selected")
		}

		links, err := dal.GetBacklinks(file.Id)
		if err != nil {
			return eris.Wrap(err, "could not get backlinks")
		}

		link, err := selectinput.Run("Select backlink", filesToChoices(links))
		if err != nil {
			return eris.Wrap(err, "failed to select backlink")
		}
		subAction, err := selectinput.Run("Select action", []selectinput.Choice{
			{Title: "Edit", Value: "edit"},
			{Title: "View", Value: "view"},
		})
		if err != nil {
			return eris.Wrap(err, "failed to select sub action")
		}

		if subAction.Value == "edit" {
			editor := utils.GetEditor()
			cmd := exec.Command(editor, link.Value)

			err = cmd.Start()
			if err != nil {
				return eris.Wrap(err, "could not open file in editor")
			}

			return nil
		}

		if subAction.Value == "view" {

			data, err := ioutil.ReadFile(link.Value)
			if err != nil {
				return eris.Wrap(err, "could not read file")
			}
			metadata := models.Fm{}
			err = frontmatter.Unmarshal(data, &metadata)
			if err != nil {
				return eris.Wrap(err, "could not unmarkshal frontmatter")
			}

			r, _ := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
			)

			out, err := r.Render(fmt.Sprintf("# %s\n%s", metadata.Title, metadata.Content))
			if err != nil {
				return eris.Wrap(err, "failed to render markdown file")
			}
			fmt.Print(out)
			return nil
		}
	}

	/*
	 */

	return nil
}

func filesToChoices(input []models.File) []selectinput.Choice {
	paths := make([]selectinput.Choice, 0)

	for _, r := range input {
		paths = append(paths, selectinput.Choice{Title: r.Path, Value: r.Path})
	}

	return paths
}