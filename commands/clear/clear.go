// Package clear removes stuff.
package clear

import (
	"fmt"
	"os"

	"github.com/hjertnes/roam/commands/help"
	iinit "github.com/hjertnes/roam/commands/init"
	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/state"
	"github.com/rotisserie/eris"
)

// Clear is the exposed type.
type Clear struct {
	state *state.State
}

// Run runs the main command and figures out what to do.
func (c *Clear) Run() error {
	if len(c.state.Arguments) == constants.Two {
		help.Run([]string{})

		return nil
	}

	switch c.state.Arguments[2] {
	case "config":
		err := os.RemoveAll(fmt.Sprintf("%s/.config", c.state.Path))
		if err != nil {
			return eris.Wrap(err, "failed to remove config")
		}

		err = iinit.Run(c.state.Path)
		if err != nil {
			return eris.Wrap(err, "failed to re-create config")
		}

		return nil
	case "database":
		err := c.state.Dal.Clear()
		if err != nil {
			return eris.Wrap(err, "failed to clear database")
		}

		return nil
	default:
		help.Run([]string{})

		return nil
	}

	return nil
}

// New is the constructor.
func New(path string, args []string) (*Clear, error) {
	s, err := state.New(path, args)
	if err != nil {
		return nil, eris.Wrap(err, "could not create state")
	}

	return &Clear{
		state: s,
	}, nil
}
