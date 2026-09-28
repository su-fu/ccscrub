package discover

import (
	"io/fs"
	"os"
	"path/filepath"
)

func Find(root string) ([]string, error) {
	var result []string
	_, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
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
		return nil, err
	}

	return result, nil
}
