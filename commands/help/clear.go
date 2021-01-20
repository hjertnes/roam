package help

import "fmt"

func clear() {
	fmt.Println("Removes setup")
	fmt.Println()
	fmt.Println("usage:")
	fmt.Println("  roam clear database\t removes database cache")
	fmt.Println("  roam clear config\t removes config directory")
	fmt.Println()
}
