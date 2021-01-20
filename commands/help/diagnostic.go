package help

import "fmt"

func diagnostic() {
	fmt.Println("This checks your roam for files with problems")
	fmt.Println("It checks that all the front matter is valid")
	fmt.Println("It checks that all the links resolve to a single file")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  roam diagnostic")
	fmt.Println()
}
