// Package find finds stuff.
package find

import (
	"fmt"
	"os"
	"strings"

	"github.com/hjertnes/roam/commands/help"
	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	"github.com/hjertnes/roam/widgets/selectinput"
	spinner2 "github.com/hjertnes/roam/widgets/spinner"
	"github.com/hjertnes/roam/widgets/textinput"
	"github.com/rotisserie/eris"
)

// Find is the exported type.
type Find struct {
	state         *state.State
	subcommand    string
	backlinksFlag bool
	linksFlag     bool
}

const query = "query"

const edit = "edit"

const view = "view"

// New is the constructor.
func New(path string) (*Find, error) {
	s, err := state.New(path)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to create state")
	}

	f := &Find{
		state:         s,
		backlinksFlag: false,
		linksFlag:     false,
	}

	if len(os.Args) == 2 || strings.HasPrefix(os.Args[2], "--") {
		f.subcommand = query
	} else {
		subCommand := os.Args[2]

		if subCommand == query || subCommand == edit || subCommand == view {
			f.subcommand = subCommand
		} else {
			help.Run()
			os.Exit(0)
		}
	}

	for _, arg := range os.Args {
		if arg == "--backlinks" {
			f.backlinksFlag = true
		}

		if arg == "--links" {
			f.linksFlag = true
		}
	}

	return f, nil
}

func (f *Find) getResults(query string) ([]models.File, error) {
	result, err := f.state.Dal.FindFileFuzzy(query)
	if err != nil {
		return make([]models.File, 0), eris.Wrap(err, "failed to search for files in database")
	}

	return result, nil
}

func (f *Find) query(files []models.File) {
	for i := range files {
		fmt.Printf("- %s\n", files[i].Path)
	}
}

func (f *Find) selectFile(result []models.File) (*models.File, error) {
	var choice *models.File

	c, err := selectinput.Run(utils.FilesToChoices(result), "Select file")
	if err != nil {
		return nil, eris.Wrap(err, "failed to get selection from user")
	}

	for i := range result {
		if result[i].ID == c.Value {
			choice = &result[i]

			break
		}
	}

	if choice == nil {
		return nil, eris.Wrap(errs.ErrNoValue, "this should not happen")
	}

	return choice, nil
}

// Run is the entrypoint.
func (f *Find) Run() error {
	search, err := textinput.Run("Search for file")
	if err != nil {
		return eris.Wrap(err, "failed to get search input from user")
	}

	spinner, err := spinner2.Run("")
	if err != nil {
		return eris.Wrap(err, "failed to build spinner")
	}

	err = spinner.Start()
	if err != nil {
		return eris.Wrap(err, "failed to start spinner")
	}

	result, err := f.getResults(search)
	if err != nil {
		return eris.Wrap(err, "failed to get search result")
	}

	err = spinner.Stop()
	if err != nil {
		return eris.Wrap(err, "failed to stop spinner")
	}

	if f.subcommand == query && !f.backlinksFlag && !f.linksFlag {
		f.query(result)

		return nil
	}

	var choice *models.File

	if len(result) == constants.Zero {
		fmt.Println("No files matches your query")
	} else if len(result) == 1 {
		choice = &result[0]
	} else {
		choice, err = f.selectFile(result)
		if err != nil {
			return eris.Wrap(err, "failed to select file")
		}

		var links []models.File

		err = spinner.Start()
		if err != nil {
			return eris.Wrap(err, "failed to start spinner")
		}

		if f.linksFlag {
			links, err = f.state.Dal.GetLinks(choice.ID)
			if err != nil {
				return eris.Wrap(err, "failed to get links")
			}
		}

		if f.backlinksFlag {
			links, err = f.state.Dal.GetBacklinks(choice.ID)
			if err != nil {
				return eris.Wrap(err, "failed to get links")
			}
		}

		err = spinner.Stop()
		if err != nil {
			return eris.Wrap(err, "failed to stop spinner")
		}

		if f.subcommand == query {
			f.query(links)

			return nil
		}

		if (f.linksFlag || f.backlinksFlag) && (f.subcommand == edit || f.subcommand == view) {
			switch len(links) {
			case 0:
				fmt.Println("No files match your query")
			case 1:
				choice = &links[0]
			default:
				choice, err = f.selectFile(links)
				if err != nil {
					return eris.Wrap(err, "failed to select file")
				}
			}
		}
	}

	if f.subcommand == edit {
		err = utils.EditNote(choice.Path)
		if err != nil {
			return eris.Wrap(err, "failed to edit file")
		}

		return nil
	}

	if f.subcommand == view {
		err = utils.ViewNote(choice.Path)
		if err != nil {
			return eris.Wrap(err, "failed to edit file")
		}

		return nil
	}

	return nil
}
