// Package configedit lets you configedit config files
package configedit

import (
	"fmt"
	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
)

// Run is the entry point.
func Run(path string) error {
	configFile := fmt.Sprintf("%s/.config/config.yaml", path)
	err := utils.EditFile(configFile)
	if err != nil {
		return eris.Wrap(err, "could not open config in editor")
	}

	return nil
}
