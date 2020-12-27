package utils

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/widgets/selectinput"
	"github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
	"io/ioutil"
	"os"
	"os/exec"
)

// GetPath returns the value of the ROAM enviornment variable or a default value if not set
func GetPath() string {
	path, isSet := os.LookupEnv("ROAM")

	if !isSet {
		return utils.ExpandTilde("~/txt/roam2")
	}

	return utils.ExpandTilde(path)
}

// GetEditor returns the value of the EDITOR enviornment variable or a default value if not set
func GetEditor() string {
	editor, isSet := os.LookupEnv("EDITOR")

	if !isSet {
		return "emacs"
	}

	return editor
}

func FilesToChoices(input []models.File) []selectinput.Choice {
	paths := make([]selectinput.Choice, 0)

	for _, r := range input {
		paths = append(paths, selectinput.Choice{Title: r.Path, Value: r.Id})
	}

	return paths
}

func EditNote(path string) error{
	editor := GetEditor()
	cmd := exec.Command(editor, path)

	err := cmd.Start()
	if err != nil {
		return eris.Wrap(err, "could not open file in editor")
	}

	return nil
}
func ViewNote(path string) error{
	data, err := ioutil.ReadFile(path)
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

func GetFile(files []models.File, id string) (*models.File, error){
	var file *models.File
	for _, r := range files {
		if r.Id == id {
			file = &r
			break
		}
	}

	if file == nil{
		eris.Wrap(errs.NotFound, "no match")
	}

	return file, nil
}