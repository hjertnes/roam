package commands

import (
	"fmt"
	"github.com/hjertnes/roam/utils"
	"os"
	"os/exec"
)




func Edit(path string){
	if len(os.Args) == 2 {
		Help()
	}

	switch os.Args[2] {
	case "config":
		openConfig(path)
		return
	default:
		Help()
		return
	}
}

func openConfig(path string) {
	editor := utils.GetEditor()
	configFile := fmt.Sprintf("%s/.config/config.yaml", path)
	cmd := exec.Command(editor, configFile)

	err := cmd.Start()
	utils.ErrorHandler(err)
}
