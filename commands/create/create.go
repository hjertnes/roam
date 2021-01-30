// Package create creates stuff.
package create

import (
	"fmt"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/utils/pathutils"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/hjertnes/roam/commands/help"
	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	"github.com/hjertnes/roam/widgets/selectinput"
	"github.com/hjertnes/roam/widgets/textinput"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
)

// Create is the exported type.
type Create struct {
	state *state.State
}

// Run runs the command and figures out what to do.
func (c *Create) RunDaily() error {
	switch len(c.state.Arguments) {
	case constants.Two:
		return c.dailyToday()
	case constants.Three:
		return c.daily(c.state.Arguments[2])
	default:
		help.Run([]string{})

		return nil
	}
}

// CreateFile creates a new note.
func (c *Create) CreateFile() error {
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

	err = c.createFile(filepath, title, templatedata)
	if err != nil {
		return eris.Wrap(err, "failed to create file")
	}

	err = utils.EditFile(fmt.Sprintf("%s/%s", c.state.Path, filepath))
	if err != nil {
		return eris.Wrap(err, "could not open file in EDITOR")
	}

	return nil
}

// New is the constructor.
func New(path string, args []string) (*Create, error) {
	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "could not create state")
	}

	return &Create{
		state: s,
	}, nil
}

func (c *Create) dailyToday() error {
	return c.daily(time.Now().Format(c.state.Conf.DateFormat))
}

// RunImport runs a import.
func (c *Create) RunImport() error {
	var file string

	dryRun := false

	for i := range c.state.Arguments {
		if i <= 1 {
			continue
		}

		if c.state.Arguments[i] == "--dry" {
			dryRun = true
		} else {
			file = c.state.Arguments[i]
		}
	}

	if file == "" {
		help.Run([]string{})
		return nil
	}

	err := c.doImport(file, dryRun)
	if err != nil {
		return eris.Wrap(err, "failed to import")
	}

	return nil
}

func (c *Create) doImport(file string, dryRun bool) error {
	importFile, err := ioutil.ReadFile(file) // #nosec
	if err != nil {
		return eris.Wrap(err, "Failed to read import file")
	}

	importContent := strings.Split(string(importFile), "\n")

	fileContent := make([]string, 0)
	sepCounter := 0

	counter := 0

	for _, line := range importContent {
		if line == "---" {
			sepCounter++
		}

		if sepCounter == constants.Three {
			err = c.writeImport(fileContent, dryRun)
			if err != nil {
				return eris.Wrap(err, "failed to import")
			}

			fileContent = make([]string, 0)

			sepCounter = 1

			counter++
		}

		fileContent = append(fileContent, line)
	}

	err = c.writeImport(fileContent, dryRun)
	if err != nil {
		return eris.Wrap(err, "failed to import")
	}

	counter++

	fmt.Printf("Imported %v notes\n", counter)

	return nil
}

func readfileImport(data string) (*models.ImportFrontmatter, error) {
	metadata := models.ImportFrontmatter{}

	err := frontmatter.Unmarshal([]byte(data), &metadata)
	if err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal frontmatter")
	}

	return &metadata, nil
}

func (c *Create) writeImport(fileContent []string, dryRun bool) error {
	data := strings.Join(fileContent, "\n")

	metadata, err := readfileImport(data)
	if err != nil {
		return eris.Wrap(err, "failed to read file for import")
	}

	exist, err := c.state.Dal.FileExists(metadata.Path)
	if err != nil {
		return eris.Wrap(err, "failed to check if file xist")
	}

	if exist {
		return eris.Wrapf(errs.ErrDuplicate, "file exist %s", metadata.Path)
	}

	parent := pathutils.New(fmt.Sprintf("%s/%s", c.state.Path, metadata.Path)).GetParent()

	if !dryRun {
		err = os.MkdirAll(parent, constants.FolderPermission)
		if err != nil {
			return eris.Wrap(err, "failed to create parent dir")
		}
	}

	p := "false"
	if metadata.Private {
		p = "true"
	}

	d := []string{
		"---",
		fmt.Sprintf(`title: "%s"`, metadata.Title),
		fmt.Sprintf(`private: %s`, p),
		"---",
		metadata.Content,
	}

	if !dryRun {
		err = c.createFile(metadata.Path, metadata.Title, []byte(strings.Join(d, "\n")))
		if err != nil {
			return eris.Wrap(err, "failed to write file for import")
		}
	} else {
		if utilslib.FileExist(fmt.Sprintf("%s/%s", c.state.Path, metadata.Path)) {
			fmt.Printf("Filename %s exist\n", metadata.Path)
		}

		if !strings.HasSuffix(metadata.Path, ".md"){
			fmt.Println("Path doesn't end in .md\n")
		}

		if strings.HasSuffix(metadata.Path, "/.md"){
			fmt.Println("Path ends in /.md seems like you forgot a filename\n")
		}
	}

	return nil
}

func (c *Create) daily(date string) error {
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

	err := utils.EditFile(fullFilename)
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

	err := ioutil.WriteFile(filepath, []byte(noteText), constants.FilePermission)
	if err != nil {
		return eris.Wrap(err, "failed to write file")
	}

	err = c.state.Dal.CreateFile(filepath, title, noteText, false)
	if err != nil {
		return eris.Wrap(err, "failed to create file in database")
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