// Package report prints reports.
package report

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	spinner2 "github.com/hjertnes/roam/widgets/spinner"
	"github.com/rotisserie/eris"
)

// Run lists files and their links.
func Run(path string) error {
	s, err := state.New(path)
	if err != nil {
		return eris.Wrap(err, "Failed to create state")
	}

	spinner, err := spinner2.Run("")
	if err != nil {
		return eris.Wrap(err, "failed to create spinner")
	}

	err = spinner.Start()
	if err != nil {
		return eris.Wrap(err, "failed to start spinner")
	}

	files, err := s.Dal.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	output := make([]string, 0)

	for i := range files {
		output, err = buildReport(s, &files[i], output)

		if err != nil {
			return eris.Wrap(err, "failed to build report")
		}
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)

	out, err := r.Render(strings.Join(output, "\n"))
	if err != nil {
		return eris.Wrap(err, "failed to render markdown file")
	}

	err = spinner.Stop()
	if err != nil {
		return eris.Wrap(err, "failed to stop spinner")
	}

	fmt.Print(out)

	return nil
}

func buildReport(s *state.State, file *models.File, output []string) ([]string, error) {
	output = append(output, fmt.Sprintf("# %s", file.Path))

	links, err := s.Dal.GetLinks(file.ID)
	if err != nil {
		return output, eris.Wrap(err, "failed to get list of links")
	}

	backlinks, err := s.Dal.GetBacklinks(file.ID)
	if err != nil {
		return output, eris.Wrap(err, "failed to get list of backlinks")
	}

	output = append(output, "## Links")

	output = utils.PrintListOfLinks(output, links)

	output = append(output, "## Backlinks")

	output = utils.PrintListOfLinks(output, backlinks)

	return output, nil
}
