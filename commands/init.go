package commands

import (
	"fmt"
	"github.com/hjertnes/roam/configuration"
	"github.com/hjertnes/roam/utils"
	utilslib "github.com/hjertnes/utils"
	"os"
)

func Init(path string){
	configFolder := fmt.Sprintf("%s/.config", path)
	configFile := fmt.Sprintf("%s/config.yaml", configFolder)
	if !utilslib.FileExist(configFolder){
		err := os.Mkdir(configFolder, 0755)
		utils.ErrorHandler(err)
	}

	if !utilslib.FileExist(configFile){
		err := configuration.CreateConfigurationFile(configFile)
		utils.ErrorHandler(err)
	}
}
