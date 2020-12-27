package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
)

const two = 2

// Edit opens something in editor.
func Edit(path string) error {
	if len(os.Args) == two {
		Help()
	}

	switch os.Args[2] {
	case "config":
		editor := utils.GetEditor()
		configFile := fmt.Sprintf("%s/.config/config.yaml", path)
		cmd := exec.Command(editor, configFile) // #nosec G204

		err := cmd.Start()
		if err != nil {
			return eris.Wrap(err, "could not open config in editor")
		}

		return nil
	default:
		Help()

		return nil
	}
}
