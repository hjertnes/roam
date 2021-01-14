package diagnostic

import (
	"fmt"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	utilslib "github.com/hjertnes/utils"
	"strings"

	"github.com/rotisserie/eris"
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

	for _, file := range files {
		if !utilslib.FileExist(file.Path) {
			fmt.Printf("%s: doesn't exist\n", file.Path)

			continue
		}

		metadata, err := utils.Readfile(file.Path)
		if err != nil {
			fmt.Printf("%s could not read file\n", file.Path)
			fmt.Println("most likely invalid front matter")

			continue
		}

		links := utils.NoteLinkRegexp.FindAllString(metadata.Content, -1)

		for _, link := range links {
			clean := utils.CleanLink(link)

			if strings.HasPrefix(clean, "/") {
				exist1, err := s.Dal.FileExists(fmt.Sprintf("%s%s.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				exist2, err := s.Dal.FileExists(fmt.Sprintf("%s%s/index.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				if !exist1 && !exist2 {
					fmt.Printf("%s no matches for link %s", file.Path, clean)
				}

				continue
			}

			matches, err := s.Dal.FindFileExact(clean)
			if err != nil {
				return eris.Wrap(err, "failed to search for link")
			}

			if len(matches) == 0 {
				fmt.Printf("%s no matches for link %s\n", file.Path, clean)

				continue
			}

			if len(matches) > 1 {
				fmt.Printf("%s more than one match for link %s\n", file.Path, clean)

				continue
			}
		}
	}

	return nil
}
