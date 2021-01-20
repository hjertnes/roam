package help

import "fmt"

func report() {
	fmt.Println("This prints simple report about your roam")
	fmt.Println("It prints the title of all your files, and all its links and backlinks")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam report")
	fmt.Println("")
}
