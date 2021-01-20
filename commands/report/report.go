// Package report prints reports.
package report

import (
	"fmt"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	spinner2 "github.com/hjertnes/roam/widgets/spinner"
	"github.com/rotisserie/eris"
	"strings"
)

type Report struct {
	state *state.State
}

func New(path string, args []string) (*Report, error) {
	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "could not create state")
	}

	return &Report{
		state: s,
	}, nil
}

// Run lists files and their links.
func (r *Report) Run() error {
	spinner, err := spinner2.Run("")
	if err != nil {
		return eris.Wrap(err, "failed to create spinner")
	}

	err = spinner.Start()
	if err != nil {
		return eris.Wrap(err, "failed to start spinner")
	}

	files, err := r.state.Dal.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	output := make([]string, 0)

	for i := range files {
		output, err = r.buildReport(&files[i], output)

		if err != nil {
			return eris.Wrap(err, "failed to build report")
		}
	}

	err = utils.RenderMarkdown(strings.Join(output, "\n"))
	if err != nil{
		return eris.Wrap(err, "failed to render")
	}

	return nil
}

func (r *Report) buildReport(file *models.File, output []string) ([]string, error) {
	output = append(output, fmt.Sprintf("# %s", file.Path))

	links, err := r.state.Dal.GetLinks(file.ID, true)
	if err != nil {
		return output, eris.Wrap(err, "failed to get list of links")
	}

	backlinks, err := r.state.Dal.GetBacklinks(file.ID, true)
	if err != nil {
		return output, eris.Wrap(err, "failed to get list of backlinks")
	}

	output = append(output, "## Links")

	output = printListOfLinks(output, links)

	output = append(output, "## Backlinks")

	output = printListOfLinks(output, backlinks)

	return output, nil
}

func printListOfLinks(output []string, links []models.File) []string {
	if len(links) == 0 {
		output = append(output, "No links")
	}

	for _, link := range links {
		output = append(output, fmt.Sprintf("- <%s>\n", link.Path))
	}

	return output
}