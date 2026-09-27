package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/su-fu/ccscrub/internal/i18n"
	"os"
)

const langEnv = "CCSCRUB_LANG"

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, i18n.T(resolveLang("", os.Getenv(langEnv)), i18n.MsgUsage))
		return 2
	}

	switch args[1] {
	case "scan":
		fmt.Println("scanに来た")
		fs := flag.NewFlagSet("scan", flag.ContinueOnError)
		lang := addLangFlag(fs)
		err := fs.Parse(args[2:])
		resolvedLang := resolveLang(*lang, os.Getenv(langEnv))
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		if err != nil {
			return 2
		}
		fmt.Println(resolvedLang)
		return 0
	case "fix":
		fmt.Println("fixに来た")
		fs := flag.NewFlagSet("fix", flag.ContinueOnError)
		dryRun := fs.Bool("dry-run", false, "show what would change without modifying any files")
		lang := addLangFlag(fs)
		err := fs.Parse(args[2:])
		resolvedLang := resolveLang(*lang, os.Getenv(langEnv))
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		if err != nil {
			return 2
		}
		fmt.Println(resolvedLang)
		fmt.Println(*dryRun)
		return 0
	default:
		fmt.Fprintln(os.Stderr, i18n.T(resolveLang("", os.Getenv(langEnv)), i18n.MsgUnknownSubcommand, args[1]))
		fmt.Fprintln(os.Stderr, i18n.T(resolveLang("", os.Getenv(langEnv)), i18n.MsgUsage))
		return 2
	}
}

func addLangFlag(fs *flag.FlagSet) *string {
	return fs.String("lang", "", "message language: en, ja, zh (default en)")
}

func resolveLang(flagLang string, envLang string) string {
	resolvedLang := flagLang
	if resolvedLang == "" {
		resolvedLang = envLang
	}
	if resolvedLang == "" {
		resolvedLang = "en"
	}
	return resolvedLang
}
