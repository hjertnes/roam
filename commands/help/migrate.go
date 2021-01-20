package help

import "fmt"

func migrate() {
	fmt.Println("This makes sure the database schema is up to date")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam migrate")
	fmt.Println("")
}