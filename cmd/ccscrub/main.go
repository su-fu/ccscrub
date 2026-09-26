package main

import (
	"errors"
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
		fs := flag.NewFlagSet("scan", flag.ContinueOnError)
		lang := addLangFlag(fs)
		err := fs.Parse(args[2:])
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		if err != nil {
			return 2
		}
		fmt.Println(*lang)
		return 0
	case "fix":
		fmt.Println("fixに来た")
		fs := flag.NewFlagSet("fix", flag.ContinueOnError)
		dryRun := fs.Bool("dry-run", false, "show what would change without modifying any files")
		lang := addLangFlag(fs)
		err := fs.Parse(args[2:])
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		if err != nil {
			return 2
		}
		fmt.Println(*lang)
		fmt.Println(*dryRun)
		return 0
	default:
		fmt.Println("defaultに来た")
		return 2
	}
}

func addLangFlag(fs *flag.FlagSet) *string {
	return fs.String("lang", "", "message language: en, ja, zh (default en)")
}
