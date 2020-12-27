package commands

import "fmt"

const version = "0.0.11-dev"

func Help() {
	fmt.Printf("roam\n")
	fmt.Printf("Version: %s\n", version)
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
	fmt.Printf("\tfind\tsearch for a note to view, edit or show the backlinks of\n")
	fmt.Printf("\tdiagnostic\tshows issues with your notes\n")
	fmt.Printf("\tedit")
	fmt.Printf("\t\tconfig\topens config file in EDITOR")
}
