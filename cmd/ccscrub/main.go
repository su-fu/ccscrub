package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/su-fu/ccscrub/internal/discover"
	"github.com/su-fu/ccscrub/internal/i18n"
	"github.com/su-fu/ccscrub/internal/jsonl"
)

const langEnv = "CCSCRUB_LANG"
const defaultMinAge = 30 * time.Minute

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
		fset := flag.NewFlagSet("scan", flag.ContinueOnError)
		lang := addLangFlag(fset)
		minAge := addMinAgeFlag(fset)
		err := fset.Parse(args[2:])
		resolvedLang := resolveLang(*lang, os.Getenv(langEnv))
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		if err != nil {
			return 2
		}
		if *minAge < 0 {
			fmt.Fprintln(os.Stderr, i18n.T(resolvedLang, i18n.MsgInvalidMinAge, *minAge))
			return 2
		}
		fmt.Println(resolvedLang)

		home, err := os.UserHomeDir()
		if err != nil {
			return 2
		}

		root := filepath.Join(home, ".claude", "projects")
		now := time.Now()
		paths, skippedList, err := discover.Find(root, *minAge, now)
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Fprintln(os.Stderr, i18n.T(resolvedLang, i18n.MsgPathNotFound, root))
			return 2
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, i18n.T(resolvedLang, i18n.MsgPathAccessError, root, err))
			return 2
		}

		for _, p := range paths {
			n, err := countLines(p)
			if err != nil {
				fmt.Fprintln(os.Stderr, i18n.T(resolvedLang, i18n.MsgCannotRead, p, err))
				continue
			}
			fmt.Printf("%s\t%d\n", p, n)
		}

		for _, skipped := range skippedList {
			switch skipped.Reason {
			case discover.SkipUnreadable:
				fmt.Fprintln(os.Stderr, i18n.T(resolvedLang, i18n.MsgCannotRead, skipped.Path, skipped.Err))
			case discover.SkipActive:
				fmt.Fprintln(os.Stderr, i18n.T(resolvedLang, i18n.MsgSkippedActive, skipped.Path, now.Sub(skipped.LastModified).Round(time.Second)))
			default:
				fmt.Fprintln(os.Stderr, i18n.T(resolvedLang, i18n.MsgSkipped, skipped.Path, skipped.Reason))
			}
		}

		return 0
	case "fix":
		fmt.Println("fixに来た")
		fset := flag.NewFlagSet("fix", flag.ContinueOnError)
		dryRun := fset.Bool("dry-run", false, "show what would change without modifying any files")
		lang := addLangFlag(fset)
		minAge := addMinAgeFlag(fset)
		err := fset.Parse(args[2:])
		resolvedLang := resolveLang(*lang, os.Getenv(langEnv))
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		if err != nil {
			return 2
		}
		if *minAge < 0 {
			fmt.Fprintln(os.Stderr, i18n.T(resolvedLang, i18n.MsgInvalidMinAge, *minAge))
			return 2
		}
		fmt.Println(resolvedLang)
		fmt.Println(*dryRun)
		fmt.Println(*minAge)
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

func addMinAgeFlag(fs *flag.FlagSet) *time.Duration {
	return fs.Duration("min-age", defaultMinAge, "skip transcripts modified within this duration, as they may be in use (e.g. 30m, 2h)")
}

func countLines(path string) (int, error) {
	count := 0
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	err = jsonl.ReadLines(f, func(num int, line []byte) error {
		count++
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
