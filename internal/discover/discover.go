package discover

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type SkipReason string

const (
	SkipUnreadable SkipReason = "unreadable"
	SkipActive     SkipReason = "active"
)

type SkippedEntry struct {
	Path         string
	Reason       SkipReason
	Err          error
	LastModified time.Time
}

func Find(root string, minAge time.Duration, now time.Time) ([]string, []SkippedEntry, error) {
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
		info, err := d.Info()
		if err != nil {
			skippedList = append(skippedList, SkippedEntry{Path: path, Reason: SkipUnreadable, Err: err})
			return nil
		}
		modTime := info.ModTime()
		if now.Sub(modTime) < minAge {
			skippedList = append(skippedList, SkippedEntry{Path: path, Reason: SkipActive, LastModified: modTime})
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
