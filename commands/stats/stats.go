package stats

import (
	"fmt"
	"github.com/hjertnes/roam/state"
	"github.com/rotisserie/eris"
)

// Stats shows statistics.
func Run(path string) error {
	s, err := state.New(path)
	if err != nil{
		return eris.Wrap(err, "Failed to create state")
	}

	all, public, private, links, err := s.Dal.Stats()
	if err != nil {
		return eris.Wrap(err, "failed to get stats")
	}

	fmt.Println("Stats")
	fmt.Printf("All: %v\n", all)
	fmt.Printf("Private: %v\n", private)
	fmt.Printf("Public: %v\n", public)
	fmt.Printf("Links: %v\n", links)

	return nil
}
