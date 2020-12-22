package commands

import (
	"fmt"
	"github.com/hjertnes/roam/configuration"
	"github.com/hjertnes/roam/utils"
	utilslib "github.com/hjertnes/utils"
	"io/ioutil"
	"os"
)

const defaultTemplate = `---
title: "$$TITLE$$"
private: false
---

`
const dailyTemplate = `---
title: "$$DATE$$"
private: false
---
`
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

	templatesDir := fmt.Sprintf("%s/templates", configFolder)
	if!utilslib.FileExist(templatesDir){
		err := os.Mkdir(templatesDir, 0700)
		utils.ErrorHandler(err)


		err = ioutil.WriteFile(fmt.Sprintf("%s/default.txt", templatesDir), []byte(defaultTemplate), 0700)
		utils.ErrorHandler(err)

		err = ioutil.WriteFile(fmt.Sprintf("%s/daily.txt", templatesDir), []byte(dailyTemplate), 0700)
		utils.ErrorHandler(err)
	}
}
