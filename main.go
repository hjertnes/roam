package main

import (
	"fmt"
	"github.com/hjertnes/roam/commands"
	"github.com/hjertnes/roam/configuration"
	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
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
		err := commands.Init(path)
		fmt.Println("Init failed")
		fmt.Println(eris.ToString(err, true))
		return
	case "edit":
		err := commands.Edit(path)
		fmt.Println("Edit failed")
		fmt.Println(eris.ToString(err, true))
		return
	case "migrate":
		err := commands.Migrate(path)
		fmt.Println("Migrate failed")
		fmt.Println(eris.ToString(err, true))
		return
	case "sync":
		err := commands.Sync(path)
		fmt.Println("Sync failed")
		fmt.Println(eris.ToString(err, true))
		return
	case "find":
		err := commands.FindEdit(path)
		fmt.Println("Find failed")
		fmt.Println(eris.ToString(err, true))
		return
	case "create":
		err := commands.Create(path, os.Args[2])
		fmt.Println("Create failed")
		fmt.Println(eris.ToString(err, true))
		return
	case "daily":
		if len(os.Args) == 2{
			conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
			fmt.Println("Daily failed")
			fmt.Println(eris.ToString(err, true))
			err = commands.Daily(path, time.Now().Format(conf.DateFormat))
			fmt.Println("Daily failed")
			fmt.Println(eris.ToString(err, true))
			return
		} else if len(os.Args) == 3{
			err := commands.Daily(path, os.Args[2])
			fmt.Println("Daily failed")
			fmt.Println(eris.ToString(err, true))
		} else {
			commands.Help()
		}
	case "view":
		err := commands.FindView(path)
		fmt.Println("View failed")
		fmt.Println(eris.ToString(err, true))
		return
	case "stats":
		err := commands.Stats(path)
		fmt.Println("Stats failed")
		fmt.Println(eris.ToString(err, true))
		return
	default:
		commands.Help()
		return
	}
}