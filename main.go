package main

import (
	"fmt"
	"github.com/hjertnes/roam/commands"
	"github.com/hjertnes/roam/configuration"
	"github.com/hjertnes/roam/errs"
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
		if err != nil {
			fmt.Println("Init failed")
			fmt.Println(eris.ToString(err, true))
			return
		}
		return
	case "diagnostic":
		err := commands.Diagnostic(path)
		if err != nil {
			fmt.Println("Diagnostic failed")
			fmt.Println(eris.ToString(err, true))
			return
		}
		return
	case "edit":
		err := commands.Edit(path)
		if err != nil {
			fmt.Println("Edit failed")
			fmt.Println(eris.ToString(err, true))
			return
		}
		return
	case "migrate":
		err := commands.Migrate(path)
		if err != nil {
			fmt.Println("Migrate failed")
			fmt.Println(eris.ToString(err, true))
			return
		}
		return
	case "sync":
		err := commands.Sync(path)
		if err != nil {
			fmt.Println("Sync failed")
			fmt.Println(eris.ToString(err, true))
			return
		}
		return
	case "find":
		err := commands.FindEdit(path)
		if err != nil {
			if eris.Is(err, errs.NotFound){
				fmt.Println("No matches to search query")
				return
			}
			fmt.Println("Find failed")
			fmt.Println(eris.ToString(err, true))
			return
		}
		return
	case "create":
		err := commands.Create(path, os.Args[2])
		if err != nil {
			fmt.Println("Create failed")
			fmt.Println(eris.ToString(err, true))
			return
		}
		return
	case "daily":
		if len(os.Args) == 2{
			conf, err := configuration.ReadConfigurationFile(fmt.Sprintf("%s/.config/config.yaml", path))
			if err != nil {
				fmt.Println("Daily failed")
				fmt.Println(eris.ToString(err, true))
				return
			}
			err = commands.Daily(path, time.Now().Format(conf.DateFormat))
			if err != nil {
				fmt.Println("Daily failed")
				fmt.Println(eris.ToString(err, true))
				return
			}
			return
		} else if len(os.Args) == 3{
			err := commands.Daily(path, os.Args[2])
			if err != nil {
				fmt.Println("Daily failed")
				fmt.Println(eris.ToString(err, true))
				return
			}
		} else {
			commands.Help()
		}
	case "stats":
		err := commands.Stats(path)
		if err != nil {
			fmt.Println("Stats failed")
			fmt.Println(eris.ToString(err, true))
			return
		}
		return
	default:
		commands.Help()
		return
	}
}