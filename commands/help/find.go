package help

import "fmt"

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
