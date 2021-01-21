package diagnostic

import (
	"fmt"
	"github.com/hjertnes/roam/models"
	"strings"

	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
)

type Diagnostic struct {
	state *state.State
}

func New(path string, args []string) (*Diagnostic, error) {
	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "could not create state")
	}

	return &Diagnostic{
		state: s,
	}, nil
}

func (d *Diagnostic) processFile(file *models.File) *models.Frontmatter{
	if !utilslib.FileExist(file.Path) {
		fmt.Printf("%s: doesn't exist\n", file.Path)

		return nil
	}

	metadata, err := utils.Readfile(file.Path)
	if err != nil {
		fmt.Printf("%s could not read file\n", file.Path)
		fmt.Println("most likely invalid front matter")

		return nil
	}

	if strings.Contains(metadata.Content, "\n\n") {
		fmt.Printf("%s contains a lot of newlines after eachother\n", file.Path)
	}

	if strings.Split(metadata.Content, "\n")[0] != "---"{
		fmt.Printf("%s the first line is not as expected", file.Path)

		return nil
	}

	return metadata
}

func (d *Diagnostic) processLink(link string, file *models.File) error {
	clean := utils.CleanLink(link)

	if strings.HasPrefix(clean, "/") {
		exist1, err := d.state.Dal.FileExists(fmt.Sprintf("%s%s.md", d.state.Path, clean))
		if err != nil {
			return eris.Wrap(err, "failed to check if link exists")
		}

		exist2, err := d.state.Dal.FileExists(fmt.Sprintf("%s%s/index.md", d.state.Path, clean))
		if err != nil {
			return eris.Wrap(err, "failed to check if link exists")
		}

		if !exist1 && !exist2 {
			fmt.Printf("%s no matches for link %s", file.Path, clean)
		}

		return nil
	}

	matches, err := d.state.Dal.FindFileExact(clean)
	if err != nil {
		return eris.Wrap(err, "failed to search for link")
	}

	if len(matches) == 0 {
		fmt.Printf("%s no matches for link %s\n", file.Path, clean)

		return nil
	}

	if len(matches) > 1 {
		fmt.Printf("%s more than one match for link %s\n", file.Path, clean)

		return nil
	}

	return nil
}

func (d *Diagnostic) Run() error {
	files, err := d.state.Dal.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	for _, file := range files {
		metadata := d.processFile(&file)
		if metadata == nil{
			continue
		}

		links := constants.NoteLinkRegexp.FindAllString(metadata.Content, -1)

		for _, link := range links {
			err = d.processLink(link, &file)
			if err != nil{
				return eris.Wrap(err, "failed to process links")
			}
		}
	}

	return nil
}
