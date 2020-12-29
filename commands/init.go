package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hjertnes/roam/configuration"
	utilslib "github.com/hjertnes/utils"
	"github.com/rotisserie/eris"
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

const publishTemplate = `<!DOCTYPE html>
<html lang="en-us">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>hjertnes.wiki</title>
	<link rel="stylesheet" href="/style.css">
</head>
<body>
    <header>
      <h1><a href="https://hjertnes.wiki/">hjertnes.wiki </a></h1>
    </a>
    </header>
	<main>
<article>
<h1>$$TITLE$$</h1>
<div>$$TEXT$$</div>
<div<$$BACKLINKS$$</div>

</article>
</main>
</body>
</html>`

// Init creates a config file.
func Init(path string) error {
	configFolder := fmt.Sprintf("%s/.config", path)
	configFile := fmt.Sprintf("%s/config.yaml", configFolder)

	if !utilslib.FileExist(configFolder) {
		err := os.Mkdir(configFolder, 0600)
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

	publishDir := fmt.Sprintf("%s/publish", configFolder)
	if !utilslib.FileExist(publishDir){
		err := os.Mkdir(publishDir, os.ModePerm)
		if err != nil {
			return eris.Wrap(err, "could not create publish folder")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/template.html", publishDir), []byte(publishTemplate), os.ModePerm)
		if err != nil {
			return eris.Wrap(err, "failed to write template")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/style.css", publishDir), []byte(""), os.ModePerm)
		if err != nil {
			return eris.Wrap(err, "failed to write template")
		}
	}

	templatesDir := fmt.Sprintf("%s/templates", configFolder)
	if !utilslib.FileExist(templatesDir) {
		err := os.Mkdir(templatesDir, os.ModePerm)
		if err != nil {
			return eris.Wrap(err, "could not create templates folder")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/default.txt", templatesDir), []byte(defaultTemplate), os.ModePerm)
		if err != nil {
			return eris.Wrap(err, "could not create default template")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/daily.txt", templatesDir), []byte(dailyTemplate), os.ModePerm)
		if err != nil {
			return eris.Wrap(err, "could not create daily note template")
		}
	}

	return nil
}
