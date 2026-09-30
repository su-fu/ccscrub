package discover

import (
	"io/fs"
	"os"
	"path/filepath"
)

type SkipReason string

const (
	SkipUnreadable SkipReason = "unreadable"
)

type SkippedEntry struct {
	Path   string
	Reason SkipReason
	Err    error
}

func Find(root string) ([]string, []SkippedEntry, error) {
	var result []string
	var skippedList []SkippedEntry
	_, err := os.Stat(root)
	if err != nil {
		return nil, nil, err
	}

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			skippedList = append(skippedList, SkippedEntry{Path: path, Reason: SkipUnreadable, Err: err})
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".jsonl" {
			return nil
		}
		result = append(result, path)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return result, skippedList, nil
}
