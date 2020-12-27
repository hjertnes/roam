package commands

import (
	"fmt"
	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
	"os"
	"os/exec"

)

func Edit(path string) error {
	if len(os.Args) == 2 {
		Help()
	}

	switch os.Args[2] {
	case "config":
		editor := utils.GetEditor()
		configFile := fmt.Sprintf("%s/.config/config.yaml", path)
		cmd := exec.Command(editor, configFile)
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
