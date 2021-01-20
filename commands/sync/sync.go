package sync

import (
	"fmt"
	"github.com/hjertnes/roam/models"
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

type Sync struct {
	state *state.State
}

func New(path string, args []string) (*Sync, error) {
	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "could not create state")
	}

	return &Sync{
		state: s,
	}, nil
}

func (s *Sync) Run() error{
	err := s.run()
	if err != nil{
		err = s.state.Dal.AddLog(true)
		if err != nil{
			return eris.Wrap(err, "failed to log failure")
		}
		return eris.Wrap(err, "failed to sync")
	}

	err = s.state.Dal.AddLog(false)
	if err != nil{
		return eris.Wrap(err, "failed to log success")
	}

	return nil
}

func (s *Sync) walkDir() error{
	err := filepath.Walk(s.state.Path, func(path string, info os.FileInfo, errr error) error {
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
		exist, err := s.state.Dal.FileExists(path)
		if err != nil {
			return eris.Wrap(err, "failed to check if file exists in database")
		}

		if exist {
			err = s.state.Dal.UpdateFile(path, metadata.Title, metadata.Content, metadata.Private)
			if err != nil {
				return eris.Wrap(err, "failed to update in database")
			}
		} else {
			err = s.state.Dal.CreateFile(path, metadata.Title, metadata.Content, metadata.Private)
			if err != nil {
				return eris.Wrap(err, "failed to create in database")
			}
		}

		return nil
	})

	if err != nil {
		return eris.Wrap(err, "failed to process files")
	}

	return nil
}

func (s *Sync) processLink(file *models.File, link string) (*models.File, error){
	clean := utils.CleanLink(link)

	filename := clean

	if strings.HasPrefix(clean, "/") {
		exist1, err := s.state.Dal.FileExists(fmt.Sprintf("%s%s.md", s.state.Path, clean))
		if err != nil {
			return nil, eris.Wrap(err, "failed to check if link exists")
		}

		exist2, err := s.state.Dal.FileExists(fmt.Sprintf("%s%s/index.md", s.state.Path, clean))
		if err != nil {
			return nil, eris.Wrap(err, "failed to check if link exists")
		}

		if exist1 {
			filename = fmt.Sprintf("%s%s.md", s.state.Path, clean)
		} else if exist2 {
			filename = fmt.Sprintf("%s%s/index.md", s.state.Path, clean)
		} else {
			return nil, eris.Wrap(errs.ErrNotFound, "not found")
		}
	}

	matches, err := s.state.Dal.FindFileExact(filename)
	if err != nil {
		return nil, eris.Wrap(err, "failed to search for link")
	}

	if len(matches) == 0 {
		return nil, eris.Wrap(errs.ErrNotFound, "no matches for selected link")
	}

	if len(matches) > 1 {
		return nil, eris.Wrap(errs.ErrNotFound, "too many matches for selected link")
	}

	err = s.state.Dal.AddLink(file.ID, matches[0].ID)
	if err != nil {
		return nil, eris.Wrap(err, "failed to add link")
	}

	return &matches[0], nil
}

func (s *Sync) processFile(file *models.File) error{
	if !utilslib.FileExist(file.Path) {
		return eris.Wrap(errs.ErrNotFound, "faile not found")
	}

	metadata, err := utils.Readfile(file.Path)
	if err != nil {
		return eris.Wrap(err, "could not read file")
	}

	links := constants.NoteLinkRegexp.FindAllString(metadata.Content, -1)

	currentInDatabaseLinks, err := s.state.Dal.GetLinks(file.ID, true)
	if err != nil {
		return eris.Wrap(err, "failed to get current links")
	}

	currentLinks := make([]string, 0)

	for _, link := range links {
		match, err := s.processLink(file, link)
		if err != nil {
			return eris.Wrap(err, "failed to process link")
		}
		currentLinks = append(currentLinks, match.ID)
	}

	for _, l := range currentInDatabaseLinks {
		if !contains(l.ID, currentLinks) {
			err := s.state.Dal.DeleteLink(file.ID, l.ID)
			if err != nil {
				return eris.Wrap(err, "failed to delete link")
			}
		}
	}

	return nil
}

func (s *Sync) processFiles() error{
	files, err := s.state.Dal.GetFiles()
	if err != nil {
		return eris.Wrap(err, "failed to get list of files")
	}

	for _, file := range files {
		err = s.processFile(&file)
		if err != nil{
			return eris.Wrap(err, "failed to process file")
		}
	}

	return nil
}

func (s *Sync) run() error{
	spinner, err := spinner2.Run("")
	if err != nil {
		return eris.Wrap(err, "failed to create spinner")
	}

	err = spinner.Start()
	if err != nil {
		return eris.Wrap(err, "failed to start spinner")
	}

	err = s.state.Dal.DeleteFiles()
	if err != nil {
		return eris.Wrap(err, "failed to delete files that don't exist from database")
	}

	err = s.walkDir()
	if err != nil{
		return eris.Wrap(err, "failed to walk dir")
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