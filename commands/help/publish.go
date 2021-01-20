package help

import "fmt"

func publish() {
	fmt.Println("This builds a website from your roam")
	fmt.Println("All notes with private: true are exluded")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  roam publish [options]\t\t Will write to ./output")
	fmt.Println("  roam publish [output-path] [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  --include-private\tincludes private notes")
	fmt.Println()
}
