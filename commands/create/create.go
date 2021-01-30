// Package create creates stuff.
package create

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/state"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
)

// Create is the exported type.
type Create struct {
	state *state.State
}

// New is the constructor.
func New(s *state.State) *Create {
	return &Create{
		state: s,
	}
}

// Run runs the command and figures out what to do.
/*





*/

func (c *Create) CreateFile(fp, title string, templatedata []byte) error {
	filepath := fmt.Sprintf("%s/%s", c.state.Path, fp)
	now := time.Now()

	if utilslib.FileExist(filepath) {
		return eris.Wrap(errs.ErrDuplicate, "we never overwrite")
	}

	noteText := strings.ReplaceAll(string(templatedata), "$$TITLE$$", title)
	noteText = strings.ReplaceAll(noteText, "$$DATE$$", now.Format(c.state.Conf.DateFormat))
	noteText = strings.ReplaceAll(noteText, "$$TIME$$", now.Format(c.state.Conf.TimeFormat))
	noteText = strings.ReplaceAll(noteText, "$$DATETIME$$", now.Format(c.state.Conf.DateTimeFormat))

	err := ioutil.WriteFile(filepath, []byte(noteText), constants.FilePermission)
	if err != nil {
		return eris.Wrap(err, "failed to write file")
	}

	err = c.state.Dal.CreateFile(filepath, title, noteText, false)
	if err != nil {
		return eris.Wrap(err, "failed to create file in database")
	}

	return nil
}
