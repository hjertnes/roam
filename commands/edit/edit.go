package edit

import (
	"fmt"
	"github.com/hjertnes/roam/commands/help"
	"github.com/hjertnes/roam/constants"
	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
	"os"
	"os/exec"
)

func Run(path string) error {
	if len(os.Args) == constants.Two {
		help.Run()
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
		help.Run()

		return nil
	}
}