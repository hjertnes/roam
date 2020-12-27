// Package utils contains various methods I don't have a better place for
package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/glamour"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/widgets/selectinput"
	"github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
)

// GetPath returns the value of the ROAM environment variable or a default value if not set.
func GetPath() string {
	path, isSet := os.LookupEnv("ROAM")

	if !isSet {
		return utils.ExpandTilde("~/txt/roam2")
	}

	return utils.ExpandTilde(path)
}

// GetEditor returns the value of the EDITOR enlivenment variable or a default value if not set.
func GetEditor() string {
	editor, isSet := os.LookupEnv("EDITOR")

	if !isSet {
		return "emacs"
	}

	return editor
}

// FilesToChoices maps a []models.File to []selectinput.Choice.
func FilesToChoices(input []models.File) []selectinput.Choice {
	paths := make([]selectinput.Choice, 0)

	for _, r := range input {
		paths = append(paths, selectinput.Choice{Title: r.Path, Value: r.ID})
	}

	return paths
}

// EditNote opens the specified file in EDITOR.
func EditNote(path string) error {
	editor := GetEditor()
	cmd := exec.Command(editor, path) // #nosec G204

	err := cmd.Start()
	if err != nil {
		return eris.Wrap(err, "could not open file in editor")
	}

	return nil
}

// ViewNote renders the specified note as markdown in terminal.
func ViewNote(path string) error {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return eris.Wrap(err, "could not read file")
	}

	metadata := models.Frontmatter{}

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

// GetFile returns a file with a id form a list of files.
func GetFile(files []models.File, id string) (*models.File, error) {
	var f *models.File

	for i := range files {
		if files[i].ID == id {
			f = &files[i]

			break
		}
	}

	if f == nil {
		return nil, eris.Wrap(errs.ErrNotFound, "no match")
	}

	return f, nil
}

// ConvertTemplateFiles convert TemplateFiles to Choice.
func ConvertTemplateFiles(templates []models.TemplateFile) []selectinput.Choice {
	result := make([]selectinput.Choice, 0)

	for _, f := range templates {
		result = append(result, selectinput.Choice{
			Title: f.Title,
			Value: f.Filename,
		})
	}

	return result
}
