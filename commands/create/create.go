package create

import (
	"fmt"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	"github.com/hjertnes/roam/widgets/selectinput"
	"github.com/hjertnes/roam/widgets/textinput"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Create struct {
	state *state.State
}

func (c *Create) CreateFile(filepath string) error {
	title, err := textinput.Run("The title of your note", "Title: ")
	if err != nil {
		return eris.Wrap(err, "could not get title from textinput")
	}

	template, err := selectinput.Run(
		"Select template",
		utils.ConvertTemplateFiles(c.state.Conf.Templates))

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

	err = c.createFile(filepath, title, templatedata)
	if err != nil {
		return eris.Wrap(err, "failed to create file")
	}

	editor := utils.GetEditor()

	cmd := exec.Command(editor, filepath) // #nosec G204

	err = cmd.Start()
	if err != nil {
		return eris.Wrap(err, "could not open file in EDITOR")
	}

	return nil
}

func New(path string) (*Create, error){
	s, err := state.New(path)
	if err != nil{
		return nil, eris.Wrap(err, "could not create state")
	}

	return &Create{
		state: s,
	}, nil
}
func (c *Create) DailyToday() error {
	return c.Daily(time.Now().Format(c.state.Conf.DateFormat))
}

func (c *Create) Import(file string) error{
	importFile, err := ioutil.ReadFile(file)
	if err != nil{
		return eris.Wrap(err, "Failed to read import file")
	}

	importContent := strings.Split(string(importFile), "\n")

	fileContent := make([]string, 0)
	sepCounter := 0

	for _, line := range importContent{
		if line == "---" {
			sepCounter++
		}

		if sepCounter == 3{
			err = c.writeImport(fileContent)
			if err != nil{
				return eris.Wrap(err, "failed to import")
			}
			fileContent = make([]string, 0)
			sepCounter = 1
		}

		fileContent = append(fileContent, line)
	}

	err = c.writeImport(fileContent)
	if err != nil{
		return eris.Wrap(err, "failed to import")
	}

	return nil

}

func (c *Create) writeImport(fileContent []string) error {
	data := strings.Join(fileContent, "\n")

	metadata, err := utils.ReadfileImport(data)
	if err != nil{
		return eris.Wrap(err, "failed to read file for import")
	}

	exist, err := c.state.Dal.FileExists(metadata.Path)
	if err != nil{
		return eris.Wrap(err, "failed to check if file xist")
	}

	if exist{
		return eris.Wrapf(errs.ErrDuplicate, "file exist %s", metadata.Path)
	}

	parent := utils.GetParent(fmt.Sprintf("%s/%s", c.state.Path, metadata.Path))
	err = os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		return eris.Wrap(err, "failed to create parent dir")
	}

	p := "false"
	if metadata.Private{
		p = "true"
	}

	d := []string{
		"---",
		fmt.Sprintf(`title: "%s"`, metadata.Title),
		fmt.Sprintf(`private: %s`, p),
		"---",
		"",
		"",
		metadata.Content,
	}

	err = c.createFile(metadata.Path, metadata.Title, []byte(strings.Join(d, "\n")))
	if err != nil{
		return eris.Wrap(err, "failed to write file for import")
	}

	return nil
}

func (c *Create) Daily(date string) error {
	filename := fmt.Sprintf("Daily Notes/%s.md", date)
	fullFilename := fmt.Sprintf("%s/%s", c.state.Path, filename)

	if !utilslib.FileExist(fullFilename) {
		templatedata, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/templates/%s", c.state.Path, "daily.txt"))
		if err != nil {
			return eris.Wrap(err, "failed to read template")
		}

		err = c.createFile(filename, "", templatedata)
		if err != nil {
			return eris.Wrap(err, "failed to create file")
		}
	}

	editor := utils.GetEditor()

	cmd := exec.Command(editor, fullFilename) // #nosec G204

	err := cmd.Start()
	if err != nil {
		return eris.Wrap(err, "faield to editNote daily in editor")
	}

	return nil
}


func (c *Create) createFile(fp, title string, templatedata []byte) error {
	filepath := fmt.Sprintf("%s/%s", c.state.Path, fp)
	now := time.Now()

	noteText := strings.ReplaceAll(string(templatedata), "$$TITLE$$", title)
	noteText = strings.ReplaceAll(noteText, "$$DATE$$", now.Format(c.state.Conf.DateFormat))
	noteText = strings.ReplaceAll(noteText, "$$TIME$$", now.Format(c.state.Conf.TimeFormat))
	noteText = strings.ReplaceAll(noteText, "$$DATETIME$$", now.Format(c.state.Conf.DateTimeFormat))

	err := ioutil.WriteFile(filepath, []byte(noteText), 0600)
	if err != nil {
		return eris.Wrap(err, "failed to write file")
	}

	// TODO this is where the crash is
	err = c.state.Dal.CreateFile(filepath, title, noteText, false)
	if err != nil {
		return eris.Wrap(err, "failed to create file in database")
	}

	return nil
}