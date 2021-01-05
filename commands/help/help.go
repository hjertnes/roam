package help

import (
	"fmt"
	"github.com/hjertnes/roam/constants"
)

func Run() {
	fmt.Printf("roam\n")
	fmt.Printf("Version: %s\n", constants.Version)
	fmt.Printf("A command line utility that will replace my use of org-roam\n")
	fmt.Printf("\n")
	fmt.Printf("Usage:\n")
	fmt.Printf("\troam\t[command]\n")
	fmt.Printf("\n")
	fmt.Printf("Available commands:\n")
	fmt.Printf("\thelp\tprints this text\n")
	fmt.Printf("\tinit\tcreates configuration files\n")
	fmt.Printf("\tmigrate\tsets up the database\n")
	fmt.Printf("\tsync\tbuilds search index used by find and others\n")
	fmt.Printf("\tfind\tsearch for a note to viewNote, editNote or show the backlinks of\n")
	fmt.Printf("\tdiagnostic\tshows issues with your notes\n")
	fmt.Printf("\treport\tlists your notes and the links and backlinks of them\n")
	fmt.Printf("\tpublish\tbuilds a html version of your database\n")
	fmt.Printf("\t\t defaults to exclude private notes, can be enabled with --include-privaste")
	fmt.Printf("\teditNote")
	fmt.Printf("\t\tconfig\topens config file in EDITOR")
}

