package help

import "fmt"

func log(){
	fmt.Println("This command manages the statuses logged after each sync. It tells if it failed or not")
	fmt.Println()
	fmt.Println("usage:")
	fmt.Println("  list\tshows all entries")
	fmt.Println("  clear\tremoves everything in the log table")
	fmt.Println()
}
