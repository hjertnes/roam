// Package configuration deals with loading and creating config files
package configuration

import (
	"io/ioutil"
	"path"

	"github.com/hjertnes/roam/models"
	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

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

	err = ioutil.WriteFile(filename, data, 0600)

	if err != nil {
		return eris.Wrap(err, "failed to write config file")
	}

	return nil
}
