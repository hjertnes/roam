// Package stats calculates statistics
package stats

import (
	"fmt"

	"github.com/hjertnes/roam/state"
	spinner2 "github.com/hjertnes/roam/widgets/spinner"
	"github.com/rotisserie/eris"
)

// Run shows statistics.
func Run(path string, args []string) error {
	s, err := state.New(path, args)
	if err != nil {
		return eris.Wrap(err, "Failed to create state")
	}

	spinner, err := spinner2.Run("")
	if err != nil {
		return eris.Wrap(err, "failed to create spinner")
	}

	err = spinner.Start()
	if err != nil {
		return eris.Wrap(err, "failed to start spinner")
	}

	all, public, private, links, err := s.Dal.Stats()
	if err != nil {
		return eris.Wrap(err, "failed to get stats")
	}

	err = spinner.Stop()
	if err != nil {
		return eris.Wrap(err, "failed to stop spinner")
	}

	fmt.Println("Stats")
	fmt.Printf("All: %v\n", all)
	fmt.Printf("Private: %v\n", private)
	fmt.Printf("Public: %v\n", public)
	fmt.Printf("Links: %v\n", links)

	return nil
}
