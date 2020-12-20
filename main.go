package main

import (
	"github.com/hjertnes/roam/commands"
	"github.com/hjertnes/roam/utils"
	"os"
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
	default:
		commands.Help()
		return
	}
}