package commands

import (
	"fmt"
	"github.com/hjertnes/roam/configuration"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
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
func Init(path string) error{
	configFolder := fmt.Sprintf("%s/.config", path)
	configFile := fmt.Sprintf("%s/config.yaml", configFolder)
	if !utilslib.FileExist(configFolder){
		err := os.Mkdir(configFolder, 0600)
		if err != nil{
			return eris.Wrap(err, "could not create config folder")
		}
	}

	if !utilslib.FileExist(configFile){
		err := configuration.CreateConfigurationFile(configFile)
		if err != nil{
			return eris.Wrap(err, "could not create config file")
		}
	}

	templatesDir := fmt.Sprintf("%s/templates", configFolder)
	if!utilslib.FileExist(templatesDir){
		err := os.Mkdir(templatesDir, 0600)
		if err != nil{
			return eris.Wrap(err, "could not create templates folder")
		}


		err = ioutil.WriteFile(fmt.Sprintf("%s/default.txt", templatesDir), []byte(defaultTemplate), 0600)
		if err != nil{
			return eris.Wrap(err, "could not create default template")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/daily.txt", templatesDir), []byte(dailyTemplate), 0600)
		if err != nil{
			return eris.Wrap(err, "could not create daily note template")
		}
	}

	return nil
}
