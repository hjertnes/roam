// Package version shows the version number
package version

import (
	"fmt"

	"github.com/hjertnes/roam/constants"
)

// Run runs it.
func Run() {
	fmt.Printf("Version: %s\n", constants.Version)
}
