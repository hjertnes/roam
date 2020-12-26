package commands

import (
	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
	"os/exec"
)

func FindEdit(path string) error{
	file, err := getFile(path)
	if err != nil{
		return eris.Wrap(err, "could not get file")
	}

	editor := utils.GetEditor()
	cmd := exec.Command(editor, file)

	err = cmd.Start()
	if err != nil {
		return eris.Wrap(err, "could not open file in editor")
	}

	return nil
}
