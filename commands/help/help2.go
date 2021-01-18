// Package help shows help information
package help

import (
	"fmt"
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

func template(){
	fmt.Println("This command let you manage templates")
	fmt.Println()
	fmt.Println("usage:")
	fmt.Println("  list\tlists all templates")
	fmt.Println("  add\tadd a new template")
	fmt.Println("  edit\tedites an template")
	fmt.Println("  delete\tdeletes an template")
	fmt.Println()
}

func log(){
	fmt.Println("This command manages the statuses logged after each sync. It tells if it failed or not")
	fmt.Println()
	fmt.Println("usage:")
	fmt.Println("  list\tshows all entries")
	fmt.Println("  clear\tremoves everything in the log table")
	fmt.Println()
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

func main() {
	fmt.Println("roam is a command line tool kind of like https://roamresearch.com/ and https://www.orgroam.com/")
	fmt.Println("A lot of the same concepts, like links and backlinks. But cli based instead of web or emacs")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam [command]")
	fmt.Println("")
	fmt.Println("Available commands:")
	fmt.Println("  create\tCreates a new file")
	fmt.Println("  daily\tCreates or opens a daily note file")
	fmt.Println("  import\tBulk import from a single file")
	fmt.Println("  diagnostic\tChecks your note collection for errors or other problems")
	fmt.Println("  config\t\tEasy access to configuration file")
	fmt.Println("  find\t\tQuery your roam database")
	fmt.Println("  help\t\tHelp like this one")
	fmt.Println("  init\t\tSets up initial configuration")
	fmt.Println("  migrate\tMakes sure your database structure is up to date")
	fmt.Println("  publish\tBuilds a HTML website version of your roam")
	fmt.Println("  report\tPrints all your notes links and backlinks")
	fmt.Println("  stats\t\tPrints simple statistics")
	fmt.Println("  sync\t\tWrites a cache of your roam into a postgres database used by search and others")
	fmt.Println("  version\tprints current version number")
	fmt.Println("  remove\tused to remove stuff")
	fmt.Println("  log\t\tused to manage the sync log")
	fmt.Println("  template\tused to manage templates")
	fmt.Println()
	fmt.Println("Use roam help [command] for more information about a specific command")
	fmt.Println()
	fmt.Println("All paths are relative to your roam dir, unless noted as full paths")
	fmt.Println()
}

func clear() {
	fmt.Println("Removes setup")
	fmt.Println()
	fmt.Println("usage:")
	fmt.Println("  roam clear database\t removes database cache")
	fmt.Println("  roam clear config\t removes config directory")
	fmt.Println()
}

func version() {
	fmt.Println("This just prints the current version")
}

func daily() {
	fmt.Println("This makes it easier to access daily notes")
	fmt.Println("Just a short hand for find and create")
	fmt.Println("It creates (if not existing) and opens the note automatically based on the Daily Notes template")

	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  roam daily\topens today")
	fmt.Println("  roam daily [date-string]\topens the date specified")
	fmt.Println()
}

func iimport() {
	fmt.Println("This bulk imports notes from a single markdown file")
	fmt.Println("It creates notes at the specified path if it doens't already exist")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  roam import /full/path/to/file.md")
	fmt.Println("  roam import /full/path/to/file.md --dry")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --dry doesn't write anything to the file system")
	fmt.Println()
}

func create() {
	fmt.Println("This makes it easy to create a new note in your roam")
	fmt.Println("You just give it a path, and you ")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  roam create [path]")
	fmt.Println()
	fmt.Println("Path: ")
	fmt.Println("  Is a file path relative to your roam root")
	fmt.Println()
}

func diagnostic() {
	fmt.Println("This checks your roam for files with problems")
	fmt.Println("It checks that all the front matter is valid")
	fmt.Println("It checks that all the links resolve to a single file")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  roam diagnostic")
	fmt.Println()
}

func edit() {
	fmt.Println("Opens your config file in EDITOR")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  roam config")
	fmt.Println()
}

func find() {
	fmt.Println("This command let you search your database")
	fmt.Println("It can print search results, links or backlinks")
	fmt.Println("Or let you open files from a search or links / backlinks of a file from a search")
	fmt.Println("Or just render it in the terminal instead")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam find")
	fmt.Println("  roam find [sub-command]")
	fmt.Println("  roam find [option]")
	fmt.Println("  roam find [sub-command] [option]")
	fmt.Println()
	fmt.Println("Sub-commands:")
	fmt.Println("  query\t[default] Prints the files matching the search")
	fmt.Println("  configedit\tOpens selected file in EDITOR")
	fmt.Println("  view\tRenders selected file in terminal")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --links\tIt applies the sub command to links of the selected note instead of itself")
	fmt.Println("  --backlinks\tIt applies the sub command to backlinks of the selected note instead of itself")
}

func iinit() {
	fmt.Println("This makes creates the initial configuration")
	fmt.Println("It never overwrites files")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam init")
	fmt.Println("")
}

func migrate() {
	fmt.Println("This makes sure the database schema is up to date")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam migrate")
	fmt.Println("")
}

func publish() {
	fmt.Println("This builds a website from your roam")
	fmt.Println("All notes with private: true are exluded")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam publish [options]\t\t Will write to ./output")
	fmt.Println("  roam publish [output-path] [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --include-private\tincludes private notes")
	fmt.Println()
}

func report() {
	fmt.Println("This prints simple report about your roam")
	fmt.Println("It prints the title of all your files, and all its links and backlinks")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam report")
	fmt.Println("")
}

func stats() {
	fmt.Println("This prints simple statistics about your roam")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam stats")
	fmt.Println("")
}

func sync() {
	fmt.Println("This maintains a cache of your roam in the SQL database")
	fmt.Println("It makes it faster and easier to do operations like search")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam sync")
	fmt.Println("")
	fmt.Println("If you get any errors try to run diagnostics")
	fmt.Println()
}