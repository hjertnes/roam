package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	spinner2 "github.com/hjertnes/roam/widgets/spinner"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
)

func Run(path string) error {
	s, err := state.New(path)
	if err != nil {
		return eris.Wrap(err, "Failed to create state")
	}

	err = run(path, s)
	if err != nil{
		err = s.Dal.AddLog(true)
		if err != nil{
			return eris.Wrap(err, "failed to log failure")
		}
		return eris.Wrap(err, "failed to sync")
	}

	err = s.Dal.AddLog(false)
	if err != nil{
		return eris.Wrap(err, "failed to log success")
	}

	return nil
}

func run(path string, s *state.State) error{
	spinner, err := spinner2.Run("")
	if err != nil {
		return eris.Wrap(err, "failed to create spinner")
	}

	err = spinner.Start()
	if err != nil {
		return eris.Wrap(err, "failed to start spinner")
	}

	err = s.Dal.DeleteFiles()
	if err != nil {
		return eris.Wrap(err, "failed to delete files that don't exist from database")
	}

	err = filepath.Walk(path, func(path string, info os.FileInfo, errr error) error {
		if errr != nil {
			return eris.Wrap(errr, "unknown problems parsing folder")
		}

		if info.Name() == ".DS_Store" {
			return nil
		}
		if strings.Contains(path, "/.") {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		metadata, err := utils.Readfile(path)
		if err != nil {
			return eris.Wrap(err, "could not read file")
		}
		exist, err := s.Dal.FileExists(path)
		if err != nil {
			return eris.Wrap(err, "failed to check if file exists in database")
		}

		if exist {
			err = s.Dal.UpdateFile(path, metadata.Title, metadata.Content, metadata.Private)
			if err != nil {
				return eris.Wrap(err, "failed to update in database")
			}
		} else {
			err = s.Dal.CreateFile(path, metadata.Title, metadata.Content, metadata.Private)
			if err != nil {
				return eris.Wrap(err, "failed to create in database")
			}
		}

		return nil
	})

	if err != nil {
		return eris.Wrap(err, "failed to process files")
	}

	files, err := s.Dal.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	for _, file := range files {
		if !utilslib.FileExist(file.Path) {
			continue
		}

		metadata, err := utils.Readfile(file.Path)
		if err != nil {
			continue
		}

		links := constants.NoteLinkRegexp.FindAllString(metadata.Content, -1)

		currentInDatabaseLinks, err := s.Dal.GetLinks(file.ID)
		if err != nil {
			return eris.Wrap(err, "failed to get current links")
		}

		currentLinks := make([]string, 0)

		for _, link := range links {
			clean := utils.CleanLink(link)

			filename := clean

			if strings.HasPrefix(clean, "/") {
				exist1, err := s.Dal.FileExists(fmt.Sprintf("%s%s.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				exist2, err := s.Dal.FileExists(fmt.Sprintf("%s%s/index.md", path, clean))
				if err != nil {
					return eris.Wrap(err, "failed to check if link exists")
				}

				if exist1 {
					filename = fmt.Sprintf("%s%s.md", path, clean)
				} else if exist2 {
					filename = fmt.Sprintf("%s%s/index.md", path, clean)
				} else {
					return eris.Wrap(errs.ErrNotFound, "not found")
				}
			}

			matches, err := s.Dal.FindFileExact(filename)
			if err != nil {
				return eris.Wrap(err, "failed to search for link")
			}

			if len(matches) == 0 {
				continue
			}

			if len(matches) > 1 {
				continue
			}

			err = s.Dal.AddLink(file.ID, matches[0].ID)
			if err != nil {
				return eris.Wrap(err, "failed to add link")
			}

			currentLinks = append(currentLinks, matches[0].ID)
		}

		for _, l := range currentInDatabaseLinks {
			if !contains(l.ID, currentLinks) {
				err := s.Dal.DeleteLink(file.ID, l.ID)
				if err != nil {
					return eris.Wrap(err, "failed to delete link")
				}
			}
		}
	}

	if err != nil {
		return eris.Wrap(err, "failed to sync")
	}

	err = spinner.Stop()
	if err != nil {
		return eris.Wrap(err, "failed to stop spinner")
	}

	return nil
}

func contains(id string, files []string) bool {
	c := false

	for _, f := range files {
		if f == id {
			c = true

			break
		}
	}

	return c
}