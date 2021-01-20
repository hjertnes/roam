package help

import "fmt"

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
