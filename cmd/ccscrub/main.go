package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/su-fu/ccscrub/internal/detect"
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

		secrets := 0
		files := 0
		for _, p := range paths {
			found, err := scanFile(p)
			if err != nil {
				fmt.Fprintln(os.Stderr, i18n.T(resolvedLang, i18n.MsgCannotRead, p, err))
				continue
			}
			if len(found) > 0 {
				secrets += len(found)
				files++
			}
			for _, f := range found {
				fmt.Printf("%s:%d:%d: %s %s\n", p, f.Line, f.Match.Start+1, f.Match.RuleID, f.Preview)
			}
		}
		fmt.Fprintln(os.Stderr, i18n.T(resolvedLang, i18n.MsgSummary, secrets, files))
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

		if files > 0 {
			return 1
		}
		return 0
	case "fix":
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

type finding struct {
	Line    int
	Match   detect.Match
	Preview string
}

func scanFile(path string) ([]finding, error) {
	var found []finding
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = jsonl.ReadLines(f, func(num int, line []byte) error {
		for _, m := range detect.Find(line) {
			p := detect.Preview(line, m)
			found = append(found, finding{Line: num, Match: m, Preview: p})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return found, nil
}

func redactLine(line []byte) ([]byte, bool) {
	newLine := detect.Redact(line, detect.Find(line))
	if json.Valid(line) && !json.Valid(newLine) {
		return line, false
	}
	return newLine, true
}
