package main

import (
	"fmt"
	"os"
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	if len(args) < 2 {
		return 2
	}

	switch args[1] {
	case "scan":
		fmt.Println("scanに来た")
		return 0
	case "fix":
		fmt.Println("fixに来た")
		return 0
	default:
		fmt.Println("defaultに来た")
		return 2
	}
}
