package template

import (
	"fmt"
	"github.com/hjertnes/roam/commands/help"
	"github.com/hjertnes/roam/configuration"
	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	"github.com/hjertnes/roam/widgets/selectinput"
	"github.com/hjertnes/roam/widgets/textinput"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
	"io/ioutil"
	"os"
	"strings"
)

// Find is the exported type.
type Template struct {
	state         *state.State
}

// New is the constructor.
func New(path string, args []string) (*Template, error) {
	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to create state")
	}

	f := &Template{
		state: s,
	}

	return f, nil
}

func (t *Template) Run() error{
	if len(t.state.Arguments) < 3{
		help.Run([]string{})
		return nil
	}

	switch t.state.Arguments[2] {
	case "add":
	err := t.Add()
		if err != nil{
			return eris.Wrap(err, "failed to add")
		}
	case "edit":
	err := t.Edit()
		if err != nil{
			return eris.Wrap(err, "failed to edit")
		}
	case "list":
	err := t.List()
		if err != nil{
			return eris.Wrap(err, "failed to list")
		}
	case "delete":
	err := t.Delete()
		if err != nil{
			return eris.Wrap(err, "failed to delete")
		}
	default:
		help.Run([]string{})
	}

	return nil
}

func (t *Template) List() error{
	out := make([]string, 0)

	for i := range t.state.Conf.Templates{
		out = append(out, fmt.Sprintf("- %s", t.state.Conf.Templates[i].Title))
	}

	err := utils.RenderMarkdown(strings.Join(out, "\n"))
	if err != nil{
		return eris.Wrap(err,"failed to render markdown")
	}

	return nil
}

func (t *Template) Add() error{
	title, err := textinput.Run("Title")
	if err != nil{
		return eris.Wrap(err, "could not get title")
	}

	if title == ""{
		fmt.Println("Can't use empty title")
		return nil
	}

	filename := fmt.Sprintf("%s/.config/templates/%s.txt", t.state.Path, title)
	if utilslib.FileExist(filename){
		fmt.Println("File with that name already exist")
		return nil
	}

	err = ioutil.WriteFile(filename, []byte(constants.DefaultTemplate), constants.FilePermission)
	if err != nil{
		return eris.Wrap(err, "failed to create template file")
	}

	t.state.Conf.Templates = append(t.state.Conf.Templates, models.TemplateFile{
		Title: title,
		Filename: fmt.Sprintf("%s.txt", title),
	})

	err  = configuration.WriteConfigurationFile(t.state.Conf, fmt.Sprintf("%s/.config/config.yaml", t.state.Path))
	if err != nil{
		return eris.Wrap(err, "failed to write config")
	}

	err = utils.EditFile(filename)
	if err != nil {
		return eris.Wrap(err, "could not open config in editor")
	}

	return nil
}

func (t *Template) getChoice()(*models.Choice, error){
	out := make([]models.Choice, 0)

	templates := t.state.Conf.Templates

	for i := range templates {
		out = append(out, models.Choice{
			Title: t.state.Conf.Templates[i].Title,
			Value: t.state.Conf.Templates[i].Filename,
		})
	}

	choice, err := selectinput.Run(out, "Select file")
	if err != nil{
		return nil, eris.Wrap(err, "failed to select file")
	}

	return choice, nil
}

func (t *Template) Edit() error{
	choice, err := t.getChoice()
	if err != nil{
		return eris.Wrap(err, "failed to get choice")
	}

	file := fmt.Sprintf("%s/.config/templates/%s", t.state.Path, choice.Value)
	err = utils.EditFile(file)
	if err != nil {
		return eris.Wrap(err, "could not open config in editor")
	}

	return nil
}

func (t *Template) Delete() error{
	choice, err := t.getChoice()
	if err != nil{
		return eris.Wrap(err, "failed to get choice")
	}

	templates := t.state.Conf.Templates

	templates = make([]models.TemplateFile, 0)
	for _, i := range t.state.Conf.Templates {
		if i.Filename == choice.Value {
			continue
		}

		templates = append(templates, i)
	}

	err = os.Remove(fmt.Sprintf("%s/.config/templates/%s", t.state.Path, choice.Value))
	if err != nil{
		return eris.Wrap(err, "failed to remove file")
	}

	t.state.Conf.Templates = templates

	err = configuration.WriteConfigurationFile(t.state.Conf, fmt.Sprintf("%s/.config/config.yaml", t.state.Path))
	if err != nil {
		return eris.Wrap(err, "failed to write config")
	}

	return nil
}