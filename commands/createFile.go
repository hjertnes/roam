package commands

import (
	"io/ioutil"
	"strings"
	"time"

	dal2 "github.com/hjertnes/roam/dal"
	"github.com/hjertnes/roam/models"
	"github.com/rotisserie/eris"
)

func createFile(dal *dal2.Dal, filepath, title string, templatedata []byte, conf *models.Configuration) error {
	now := time.Now()

	noteText := strings.ReplaceAll(string(templatedata), "$$TITLE$$", title)
	noteText = strings.ReplaceAll(noteText, "$$DATE$$", now.Format(conf.DateFormat))
	noteText = strings.ReplaceAll(noteText, "$$TIME$$", now.Format(conf.TimeFormat))
	noteText = strings.ReplaceAll(noteText, "$$DATETIME$$", now.Format(conf.DateTimeFormat))

	err := ioutil.WriteFile(filepath, []byte(noteText), 0600)
	if err != nil {
		return eris.Wrap(err, "failed to write file")
	}

	err = dal.Create(filepath, title, "", false)
	if err != nil {
		return eris.Wrap(err, "failed to create file in database")
	}

	return nil
}
