package configuration

import (
	"github.com/hjertnes/roam/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ReadConfigurationFile(filename string) (*models.Configuration, error){
	conf := models.Configuration{}

	data, err := ioutil.ReadFile(filename)
	if err != nil{
		return nil, err
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil{
		return nil, err
	}

	return &conf, nil
}

func CreateConfigurationFile(filename string) error{
	config := &models.Configuration{
		TimeFormat: "15:04:05Z07:00",
		DateFormat: "2006-01-02",
		DateTimeFormat: "2006-01-02T15:04:05Z07:00",
		DefaultFileExtension: "md",
		Version: 0,
		Templates: []models.TemplateFile{
			{
				Filename: "default.txt",
				Title: "Default",
			},
			{
				Filename: "daily.txt",
				Title: "Daily Note",
			},
		},
	}

	data, err := yaml.Marshal(&config)

	if err != nil{
		return err
	}

	err = ioutil.WriteFile(filename, data, 0700)

	if err != nil{
		return err
	}

	return nil
}
