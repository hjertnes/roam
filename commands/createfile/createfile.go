package createfile

import (
	"fmt"
	create2 "github.com/hjertnes/roam/commands/create"
	"github.com/hjertnes/roam/commands/help"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	"github.com/hjertnes/roam/widgets/selectinput"
	"github.com/hjertnes/roam/widgets/textinput"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
	"io/ioutil"
)

type CreateFile struct {
	state *state.State
	create *create2.Create
}


// New is the constructor.
func New(path string, args []string) (*CreateFile, error) {
	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "could not create state")
	}

	c := create2.New(s)
	return &CreateFile{
		create: c,
		state: s,
	}, nil
}

// CreateFile creates a new note.
func (c *CreateFile) Run() error {
	if len(c.state.Arguments) <= 3 {
		help.Run([]string{})
		return nil
	}

	filepath := c.state.Arguments[2]

	title, err := textinput.Run("Title")
	if err != nil {
		return eris.Wrap(err, "could not get title from textinput")
	}

	template, err := selectinput.Run(
		convertTemplateFiles(c.state.Conf.Templates), "Template")
	if err != nil {
		return eris.Wrap(err, "could not get template selection from selectinput")
	}

	if utilslib.FileExist(filepath) {
		return errs.ErrDuplicate
	}

	templatedata, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/templates/%s", c.state.Path, template.Value))
	if err != nil {
		return eris.Wrap(err, "could not read template")
	}

	err = c.create.CreateFile(filepath, title, templatedata)
	if err != nil {
		return eris.Wrap(err, "failed to create file")
	}

	err = utils.EditFile(fmt.Sprintf("%s/%s", c.state.Path, filepath))
	if err != nil {
		return eris.Wrap(err, "could not open file in EDITOR")
	}

	return nil
}

func convertTemplateFiles(templates []models.TemplateFile) []models.Choice {
	result := make([]models.Choice, 0)

	for _, f := range templates {
		result = append(result, models.Choice{
			Title: f.Title,
			Value: f.Filename,
		})
	}

	return result
}