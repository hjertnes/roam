// Package init creates the config dir and its files
package init

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hjertnes/roam/configuration"
	"github.com/hjertnes/roam/constants"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
)

// Run creates a config file.
func Run(path string) error {
	configFolder := fmt.Sprintf("%s/.config", path)
	configFile := fmt.Sprintf("%s/config.yaml", configFolder)

	if !utilslib.FileExist(path) {
		err := os.MkdirAll(path, constants.FolderPermission)
		if err != nil {
			return eris.Wrap(err, "could not create data folder")
		}
	}

	if !utilslib.FileExist(configFolder) {
		err := os.Mkdir(configFolder, constants.FolderPermission)
		if err != nil {
			return eris.Wrap(err, "could not create config folder")
		}
	}

	if !utilslib.FileExist(configFile) {
		err := configuration.CreateConfigurationFile(configFile)
		if err != nil {
			return eris.Wrap(err, "could not create config file")
		}
	}

	err := createPublishDir(configFolder)
	if err != nil {
		return eris.Wrap(err, "faield to create config folder")
	}

	err = createTemplateDir(configFolder)
	if err != nil {
		return eris.Wrap(err, "faield to create templates folder")
	}

	return nil
}

func createTemplateDir(configFolder string) error {
	templatesDir := fmt.Sprintf("%s/templates", configFolder)
	if !utilslib.FileExist(templatesDir) {
		err := os.Mkdir(templatesDir, constants.FolderPermission)
		if err != nil {
			return eris.Wrap(err, "could not create templates folder")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/default.txt", templatesDir), []byte(constants.DefaultTemplate), constants.FilePermission)
		if err != nil {
			return eris.Wrap(err, "could not create default template")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/daily.txt", templatesDir), []byte(constants.DefaultTemplate), constants.FilePermission)
		if err != nil {
			return eris.Wrap(err, "could not create daily note template")
		}
	}

	return nil
}

func createPublishDir(configFolder string) error {
	publishDir := fmt.Sprintf("%s/publish", configFolder)
	if !utilslib.FileExist(publishDir) {
		err := os.Mkdir(publishDir, constants.FolderPermission)
		if err != nil {
			return eris.Wrap(err, "could not create publish folder")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/template.html", publishDir), []byte(constants.PublishTemplate), constants.FilePermission)
		if err != nil {
			return eris.Wrap(err, "failed to write template")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/style.css", publishDir), []byte(""), constants.FilePermission)
		if err != nil {
			return eris.Wrap(err, "failed to write template")
		}
	}

	return nil
}
