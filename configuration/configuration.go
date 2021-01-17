// Package configuration deals with loading and creating config files
package configuration

import (
	"github.com/hjertnes/roam/constants"
	"io/ioutil"
	"path"

	"github.com/hjertnes/roam/models"
	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

func updateConfig(conf *models.Configuration) bool{
	changed := false
	if conf.Publish == nil{
		conf.Publish = &models.Publish{
			FilesToCopy: []string{
				".js",
				".css",
				".jpeg",
				".jpg",
				".png",
				".gif",
			},
		}
		changed = true
	}

	return changed
}

// ReadConfigurationFile returns config file from filename.
func ReadConfigurationFile(filename string) (*models.Configuration, error) {
	conf := models.Configuration{}

	data, err := ioutil.ReadFile(path.Clean(filename))
	if err != nil {
		return nil, eris.Wrap(err, "failed to read config file")
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, eris.Wrap(err, "failed to unmarshal config file")
	}

	if updateConfig(&conf){
		err := WriteConfigurationFile(&conf, filename)
		if err != nil{
			return nil, eris.Wrap(err, "failed to update config")
		}
	}

	return &conf, nil
}

// CreateConfigurationFile creates a new config file.
func CreateConfigurationFile(filename string) error {
	config := &models.Configuration{
		DatabaseConnectionString: "",
		TimeFormat:               "15:04:05Z07:00",
		DateFormat:               "2006-01-02",
		DateTimeFormat:           "2006-01-02T15:04:05Z07:00",
		DefaultFileExtension:     "md",
		Version:                  0,
		Publish: &models.Publish{
			FilesToCopy: []string{
				".js",
				".css",
				".jpeg",
				".jpg",
				".png",
				".gif",
			},
		},
		Templates: []models.TemplateFile{
			{
				Filename: "default.txt",
				Title:    "Default",
			},
			{
				Filename: "daily.txt",
				Title:    "Daily Note",
			},
		},
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		return eris.Wrap(err, "failed to marshal config file")
	}

	err = ioutil.WriteFile(filename, data, constants.FilePermission)

	if err != nil {
		return eris.Wrap(err, "failed to write config file")
	}

	return nil
}


func WriteConfigurationFile(config *models.Configuration, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return eris.Wrap(err, "failed to marshal config file")
	}

	err = ioutil.WriteFile(filename, data, constants.FilePermission)

	if err != nil {
		return eris.Wrap(err, "failed to write config file")
	}

	return nil
}
