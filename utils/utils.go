// Package utils contains various methods I don't have a better place for
package utils

import (
	"fmt"
	"github.com/hjertnes/roam/errs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/models"
	"github.com/rotisserie/eris"
)

func ErrorHandler(err error) {
	if err != nil {
		if eris.Is(err, errs.ErrNotFound) {
			fmt.Println("No matches to search query")

			os.Exit(0)
		}

		if eris.Is(err, errs.ErrNoop){
			os.Exit(0)
		}

		fmt.Println("Error")

		fmt.Println(eris.ToString(err, true))

		os.Exit(0)
	}
}

// GetEditor returns the value of the EDITOR enlivenment variable or a default value if not set.
func GetEditor() string {
	editor, isSet := os.LookupEnv("EDITOR")

	if !isSet {
		return "emacs"
	}

	return editor
}

func RenderMarkdown(data string) error{
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	out, err := r.Render(data)
	if err != nil {
		return eris.Wrap(err, "failed to render markdown")
	}

	fmt.Print(out)

	return nil
}

// CleanLink removes the [[ and ]] around links.
func CleanLink(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, "[[", ""), "]]", "")
}

// Readfile reads a file into a model.
func Readfile(path string) (*models.Frontmatter, error) {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, eris.Wrap(err, "failed to read file")
	}

	metadata := models.Frontmatter{}

	err = frontmatter.Unmarshal(data, &metadata)
	if err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal frontmatter")
	}

	return &metadata, nil
}
