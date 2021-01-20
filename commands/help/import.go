package help

import "fmt"

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
