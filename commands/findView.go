package commands

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/models"
	"github.com/rotisserie/eris"
	"io/ioutil"
)

func FindView(path string) error{
	file, err := getFile(path)
	if err != nil {
		return eris.Wrap(err, "could not get file")
	}

	data, err := ioutil.ReadFile(file)
	if err != nil{
		return eris.Wrap(err, "could not read file")
	}
	metadata := models.Fm{}
	err = frontmatter.Unmarshal(data, &metadata)
	if err != nil{
		return eris.Wrap(err, "could not unmarkshal frontmatter")
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	out, err := r.Render(fmt.Sprintf("# %s\n%s", metadata.Title, metadata.Content))
	if err != nil{
		return eris.Wrap(err, "failed to render markdown file")
	}
	fmt.Print(out)
	return nil
}