package main

import (
	"flag"
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
		fs := flag.NewFlagSet("scan", flag.ExitOnError)
		lang := fs.String("lang", "", "message language: en, ja, zh (default en)")
		err := fs.Parse(args[2:])
		if err != nil {
			return 2
		}
		fmt.Println(*lang)
		return 0
	case "fix":
		fmt.Println("fixに来た")
		return 0
	default:
		fmt.Println("defaultに来た")
		return 2
	}
}
