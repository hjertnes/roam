// Package textinput contains a terminal widget for getting a line of input from a user
package textinput

import (
	"fmt"

	"github.com/rotisserie/eris"
)

// Run runs it.
func Run(label string) (string, error) {
	var input string

	fmt.Printf("%s: ", label)

	_, err := fmt.Scan(&input)
	if err != nil {
		return input, eris.Wrap(err, "could not get text input from user")
	}

	return input, nil
}
