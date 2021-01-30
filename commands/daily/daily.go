package daily

import (
	"fmt"
	create2 "github.com/hjertnes/roam/commands/create"
	"github.com/hjertnes/roam/commands/help"
	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
	"io/ioutil"
	"time"
)

type Daily struct {
	state *state.State
	create *create2.Create
}


// New is the constructor.
func New(path string, args []string) (*Daily, error) {
	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "could not create state")
	}

	c := create2.New(s)
	return &Daily{
		create: c,
		state: s,
	}, nil
}

func (d *Daily) Run() error {
	switch len(d.state.Arguments) {
	case constants.Two:
		return d.dailyToday()
	case constants.Three:
		return d.daily(d.state.Arguments[2])
	default:
		help.Run([]string{})

		return nil
	}
}

func (c *Daily) daily(date string) error {
	filename := fmt.Sprintf("Daily Notes/%s.md", date)
	fullFilename := fmt.Sprintf("%s/%s", c.state.Path, filename)

	if !utilslib.FileExist(fullFilename) {
		templatedata, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/templates/%s", c.state.Path, "daily.txt"))
		if err != nil {
			return eris.Wrap(err, "failed to read template")
		}

		err = c.create.CreateFile(filename, "", templatedata)
		if err != nil {
			return eris.Wrap(err, "failed to create file")
		}
	}

	err := utils.EditFile(fullFilename)
	if err != nil {
		return eris.Wrap(err, "faield to editNote daily in editor")
	}

	return nil
}

func (d *Daily) dailyToday() error {
	return d.daily(time.Now().Format(d.state.Conf.DateFormat))
}