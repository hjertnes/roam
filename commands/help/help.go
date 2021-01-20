// Package help shows help information
package help

import (
	"strings"
)

func getSubCommand(args []string) string {
	subCommand := ""

	for i := range args {
		if i > 1 {
			if !strings.HasPrefix(args[i], "--") {
				subCommand = args[i]

				break
			}
		}
	}

	if len(args) > 1 && args[1] != "help" {
		subCommand = ""
	}

	return subCommand
}

var subCommands = map[string]func(){
	"create":     create,
	"clear":      clear,
	"diagnostic": diagnostic,
	"configedit":       edit,
	"find":       find,
	"init":       iinit,
	"migrate":    migrate,
	"publish":    publish,
	"report":     report,
	"stats":      stats,
	"sync":       sync,
	"daily":      daily,
	"import":     iimport,
	"version":    version,
	"template": template,
	"log": log,
}

func contains(key string) bool {
	for i := range subCommands {
		if i == key {
			return true
		}
	}

	return false
}

// Run is the entry point.
func Run(args []string) {
	subCommand := getSubCommand(args)
	if !contains(subCommand) {
		main()
	} else {
		subCommands[subCommand]()
	}
}