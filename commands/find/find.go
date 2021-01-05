package find

import (
	"fmt"
	"github.com/hjertnes/roam/models"
	"github.com/hjertnes/roam/state"


	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/utils"
	"github.com/hjertnes/roam/widgets/selectinput"
	"github.com/hjertnes/roam/widgets/textinput"
	"github.com/rotisserie/eris"
)

func Run(path string) error {
	search, err := textinput.Run("Search for a note", "Search: ")
	if err != nil {
		return eris.Wrap(err, "failed to get a search string from textinput")
	}

	fmt.Println("Loading...")

	s, err := state.New(path)
	if err != nil{
		return eris.Wrap(err, "Failed to create state")
	}


	result, err := s.Dal.FindFileFuzzy(search)
	if err != nil {
		return eris.Wrap(err, "failed to search for files in database")
	}

	fmt.Println("Loading...")

	choice, err := selectinput.Run("Select match", utils.FilesToChoices(result))
	if err != nil {
		return eris.Wrap(err, "failed to get selection from selectinput")
	}

	file, err := utils.GetFile(result, choice.Value)
	if err != nil {
		return eris.Wrap(err, "should not happen")
	}

	action, err := selectinput.Run("Select action", []models.Choice{
		{Title: "Edit", Value: "editNote"},
		{Title: "View", Value: "viewNote"},
		{Title: "Backlinks", Value: "backlinks"},
	})
	if err != nil {
		return eris.Wrap(err, "could not get action")
	}

	if action.Value == "editNote" {
		err = utils.EditNote(file.Path)
		if err != nil {
			return eris.Wrap(err, "failed to edit file")
		}
	}

	if action.Value == "viewNote" {
		err = utils.ViewNote(file.Path)
		if err != nil {
			return eris.Wrap(err, "failed to show file")
		}

		return nil
	}

	if action.Value == "backlinks" {
		err := showBacklinks(s.Dal, file, result)
		if err != nil{
			return eris.Wrap(err, "backlinks failed")
		}
	}

	return nil
}

func showBacklinks(dal dal2.Dal, file *models.File, result []models.File) error{
	links, err := dal.GetBacklinks(file.ID)
	if err != nil {
		return eris.Wrap(err, "could not get backlinks")
	}

	fmt.Println("Loading...")

	link, err := selectinput.Run("Select backlink", utils.FilesToChoices(links))
	if err != nil {
		return eris.Wrap(err, "failed to select backlink")
	}

	file2, err := utils.GetFile(result, link.Value)
	if err != nil {
		return eris.Wrap(err, "should not happen")
	}

	subAction, err := selectinput.Run("Select action", []models.Choice{
		{Title: "Edit", Value: "editNote"},
		{Title: "View", Value: "viewNote"},
	})
	if err != nil {
		return eris.Wrap(err, "failed to select sub action")
	}

	if subAction.Value == "editNote" {
		err = utils.EditNote(file2.Path)
		if err != nil {
			return eris.Wrap(err, "failed to edit file")
		}
	}

	if subAction.Value == "viewNote" {
		err = utils.ViewNote(file2.Path)
		if err != nil {
			return eris.Wrap(err, "failed to show file")
		}


	}

	return nil
}