package help

import "fmt"

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
