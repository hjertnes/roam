package main

import (
	"fmt"
	"github.com/hjertnes/roam/commands"
	"github.com/hjertnes/roam/configuration"
	"github.com/hjertnes/roam/utils"
	"os"
	"time"
)

func main(){
	path := utils.GetPath()
	if len(os.Args) == 1{
		commands.Help()
		return
	}

	switch os.Args[1] {
	case "init":
		commands.Init(path)
		return
	case "edit":
		commands.Edit(path)
		return
	case "migrate":
		commands.Migrate(path)
		return
	case "sync":
		commands.Sync(path)
		return
	case "find":
		commands.FindEdit(path)
		return
	case "create":
		commands.Create(path, os.Args[2])
	case "daily":
		if len(os.Args) == 2{
			conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
			utils.ErrorHandler(err)
			commands.Daily(path, time.Now().Format(conf.DateFormat))
			return
		} else if len(os.Args) == 3{
			commands.Daily(path, os.Args[2])
		} else {
			commands.Help()
		}
	case "view":
		commands.FindView(path)
		return
	default:
		commands.Help()
		return
	}
}