package help

import "fmt"

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
