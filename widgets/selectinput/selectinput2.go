// Package selectinput contains a terminal interface like a dropdown
package selectinput

import (
	"fmt"
	"strconv"

	"github.com/hjertnes/roam/errs"
	"github.com/hjertnes/roam/models"
	"github.com/rotisserie/eris"
)

// Run is the function that runs it.
func Run(choices []models.Choice, label string) (*models.Choice, error) {
	for i, c := range choices {
		fmt.Printf("[%v]\t %s\n", i, c.Title)
	}

	status := true

	for {
		if !status {
			fmt.Println("Invalid input try again: q to quit")
		}

		var value string

		fmt.Printf("%s: ", label)

		_, err := fmt.Scanln(&value)
		if err != nil {
			return nil, eris.Wrap(err, "failed to get user input")
		}

		if value == "q" {
			return nil, eris.Wrap(errs.ErrNoValue, "no input supplied")
		}

		valueAsInt, err := strconv.Atoi(value)
		if err != nil {
			status = false

			continue
		}

		if valueAsInt < 0 || valueAsInt >= len(choices) {
			status = false

			continue
		}

		return &choices[valueAsInt], nil
	}
}
