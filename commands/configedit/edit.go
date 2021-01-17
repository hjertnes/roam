// Package configedit lets you configedit config files
package configedit

import (
	"fmt"
	"os/exec"

	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
)

// Run is the entry point.
func Run(path string) error {
	editor := utils.GetEditor()
	configFile := fmt.Sprintf("%s/.config/config.yaml", path)
	cmd := exec.Command(editor, configFile) // #nosec G204

	err := cmd.Start()
	if err != nil {
		return eris.Wrap(err, "could not open config in editor")
	}

	return nil
}
