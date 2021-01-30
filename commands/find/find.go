// Package find finds stuff.
package find

import (
	"fmt"
	"github.com/ericaro/frontmatter"
	"github.com/hjertnes/roam/constants"
	"github.com/theckman/yacspin"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hjertnes/roam/commands/help"
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
func New(path string, args []string) (*Find, error) {
	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "Failed to create state")
	}

	f := &Find{
		state:         s,
		backlinksFlag: false,
		linksFlag:     false,
	}

	if len(s.Arguments) == 2 || strings.HasPrefix(s.Arguments[2], "--") {
		f.subcommand = query
	} else {
		subCommand := s.Arguments[2]

		if subCommand == query || subCommand == edit || subCommand == view {
			f.subcommand = subCommand
		} else {
			help.Run(args)
			return nil, errs.ErrNoop
		}
	}

	for _, arg := range s.Arguments {
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

func (f *Find) query(files []models.File) error {
	result := make([]string, 0)
	for i := range files {
		result = append(result, fmt.Sprintf("- <%s>\n", files[i].Path))
	}

	err := utils.RenderMarkdown(strings.Join(result, "\n"))
	if err != nil{
		return eris.Wrap(err, "failed to render markdown")
	}

	return nil
}

func (f *Find) selectFile(result []models.File) (*models.File, error) {
	var choice *models.File

	c, err := selectinput.Run(filesToChoices(result), "Select file")
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

func (f *Find) getLinks(choice *models.File, spinner *yacspin.Spinner) ([]models.File, error){
	err := spinner.Stop()

	var links []models.File

	err = spinner.Start()
	if err != nil {
		return links, eris.Wrap(err, "failed to start spinner")
	}

	if f.linksFlag {
		links, err = f.state.Dal.GetLinks(choice.ID, true)
		if err != nil {
			return links, eris.Wrap(err, "failed to get links")
		}
	}

	if f.backlinksFlag {
		links, err = f.state.Dal.GetBacklinks(choice.ID, true)
		if err != nil {
			return links, eris.Wrap(err, "failed to get links")
		}
	}

	err = spinner.Stop()
	if err != nil {
		return links, eris.Wrap(err, "failed to stop spinner")
	}

	return links, errs.ErrNoValue
}

func (f *Find) getChoice(result []models.File, spinner *yacspin.Spinner) (*models.File, error){
	if len(result) == constants.Zero {
		return nil, errs.ErrNoValue
	} else if len(result) == 1 {
		return &result[0], nil
	} else {
		choice, err := f.selectFile(result)
		if err != nil {
			return nil, eris.Wrap(err, "failed to select file")
		}

		links, err := f.getLinks(choice, spinner)
		if err != nil{
			return nil, eris.Wrap(err, "failed to get links")
		}

		if f.subcommand == query {
			err := f.query(result)
			if err != nil{
				return nil, eris.Wrap(err, "failed to show query result")
			}

			return nil, errs.ErrNoValue
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
					return nil, eris.Wrap(err, "failed to select file")
				}
			}
		}
	}

	return nil, errs.ErrNoValue
}

// Run is the entrypoint.
func (f *Find) Run() error {
	if f.linksFlag && f.backlinksFlag{
		fmt.Println("You can't use both backlinks and links at the same time")
		return nil
	}

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
		err := f.query(result)
		if err != nil{
			return eris.Wrap(err, "failed to show query result")
		}

		return nil
	}

	choice, err := f.getChoice(result, spinner)
	if err != nil{
		if eris.Is(err, errs.ErrNoValue){
			return nil
		}

		return eris.Wrap(err, "failed to get choice")
	}

	return f.openOrEdit(choice)
}

func (f *Find) openOrEdit(choice *models.File) error{


	if f.subcommand == edit {
		err := utils.EditFile(choice.Path)
		if err != nil {
			return eris.Wrap(err, "failed to edit file")
		}

		return nil
	}

	if f.subcommand == view {
		err := viewNote(choice.Path)
		if err != nil {
			return eris.Wrap(err, "failed to view file")
		}

		return nil
	}

	return nil
}

func viewNote(path string) error {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return eris.Wrap(err, "could not read file")
	}

	metadata := models.Frontmatter{}

	err = frontmatter.Unmarshal(data, &metadata)
	if err != nil {
		return eris.Wrap(err, "could not unmarkshal frontmatter")
	}

	err = utils.RenderMarkdown(fmt.Sprintf("# %s\n%s", metadata.Title, metadata.Content))

	if err != nil {
		return eris.Wrap(err, "failed to render markdown file")
	}

	return nil
}


func filesToChoices(input []models.File) []models.Choice {
	paths := make([]models.Choice, 0)

	for _, r := range input {
		paths = append(paths, models.Choice{Title: r.Path, Value: r.ID})
	}

	return paths
}