package report

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
	"strings"
)



func Run(path string) error {
	s, err := state.New(path)
	if err != nil{
		return eris.Wrap(err, "Failed to create state")
	}


	files, err := s.Dal.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	output := make([]string, 0)

	for _, file := range files {
		output = append(output, fmt.Sprintf("# %s", file.Path))

		links, err := s.Dal.GetLinks(file.ID)
		if err != nil {
			return eris.Wrap(err, "failed to get list of links")
		}

		backlinks, err := s.Dal.GetBacklinks(file.ID)
		if err != nil {
			return eris.Wrap(err, "failed to get list of backlinks")
		}

		output = append(output, "## Links")

		output = utils.PrintListOfLinks(output, links)

		output = append(output, "## Backlinks")

		output = utils.PrintListOfLinks(output, backlinks)
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	out, err := r.Render(strings.Join(output, "\n"))
	if err != nil {
		return eris.Wrap(err, "failed to render markdown file")
	}

	fmt.Print(out)

	return nil
}
