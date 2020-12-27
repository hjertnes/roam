package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/hjertnes/roam/configuration"
	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/utils"
	"github.com/hjertnes/roam/widgets/selectinput"
	"github.com/hjertnes/roam/widgets/textinput"
	utilslib "github.com/hjertnes/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rotisserie/eris"
)

func convertTemplateFiles(templates []models.TemplateFile) []selectinput.Choice {
	result := make([]selectinput.Choice, 0)

	for _, f := range templates {
		result = append(result, selectinput.Choice{
			Title: f.Title,
			Value: f.Filename,
		})
	}

	return result
}

func Create(path, filepath string) error {
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

	title, err := textinput.Run("The title of your note", "Title: ") //noteTitle.Run()
	if err != nil {
		return eris.Wrap(err, "could not get title from textinput")
	}

	template, err := selectinput.Run("Select template", convertTemplateFiles(conf.Templates)) //selectType.Run(conf.Templates)
	if err != nil {
		return eris.Wrap(err, "could not get template selection from selectinput")
	}

	if utilslib.FileExist(filepath) {
		return errs.Duplicate
	}

	templatedata, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/templates/%s", path, template.Value))
	if err != nil {
		return eris.Wrap(err, "could not read template")
	}

	err = createFile(dal, filepath, title, templatedata, conf)
	if err != nil {
		return eris.Wrap(err, "failed to create file")
	}

	editor := utils.GetEditor()

	cmd := exec.Command(editor, filepath)

	err = cmd.Start()
	if err != nil {
		return eris.Wrap(err, "could not open file in EDITOR")
	}

	return nil
}
