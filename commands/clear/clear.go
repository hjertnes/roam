package clear

import (
	"fmt"
	"github.com/hjertnes/roam/commands/help"
	iinit "github.com/hjertnes/roam/commands/init"
	"github.com/hjertnes/roam/state"
	"github.com/rotisserie/eris"
	"os"
)

type Clear struct {
	state *state.State
}

func (c *Clear) Run() error {
	if len(os.Args) == 2{
		help.Run()
		os.Exit(0)
	}

	switch os.Args[2] {
	case "config":
		err := os.RemoveAll(fmt.Sprintf("%s/.config", c.state.Path))
		if err != nil{
			return eris.Wrap(err, "failed to remove config")
		}

		err = iinit.Run(c.state.Path)
		if err != nil{
			return eris.Wrap(err, "failed to re-create config")
		}

		os.Exit(0)
	case "database":
		err := c.state.Dal.Clear()
		if err != nil{
			return eris.Wrap(err, "failed to clera database")
		}
		os.Exit(0)
	default:
		help.Run()
		os.Exit(0)
	}

	return nil
}


func New(path string) (*Clear, error){
	s, err := state.New(path)
	if err != nil{
		return nil, eris.Wrap(err, "could not create state")
	}

	return &Clear{
		state: s,
	}, nil
}
