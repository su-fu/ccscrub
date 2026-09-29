package discover

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestFind(t *testing.T) {
	root := t.TempDir()

	paths := []string{
		"projA/aaa.jsonl",
		"projA/bbb.jsonl",
		"projA/memory/note.md",
		"projB/ccc.jsonl",
		"projB/ccc/tool-results/x.txt",
	}

	for _, p := range paths {
		full := filepath.Join(root, p)
		err := os.MkdirAll(filepath.Dir(full), 0o755)
		if err != nil {
			t.Fatalf("MkdirAll(%q): %v", filepath.Dir(full), err)
		}
		err = os.WriteFile(full, nil, 0o644)
		if err != nil {
			t.Fatalf("WriteFile(%q): %v", full, err)
		}
	}

	want := []string{
		filepath.Join(root, "projA", "aaa.jsonl"),
		filepath.Join(root, "projA", "bbb.jsonl"),
		filepath.Join(root, "projB", "ccc.jsonl"),
	}

	got, err := Find(root)
	if err != nil {
		t.Fatalf("Find(%q) error = %v", root, err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("Find(%q) = %v, want %v", root, got, want)
	}
}

func TestFindMissingRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "nope")
	paths, err := Find(root)
	if err == nil {
		t.Fatalf("Find(%q) error = %v, want an error", root, err)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Find(%q) error = %v, want fs.ErrNotExist", root, err)
	}
	if paths != nil {
		t.Errorf("Find(%q) paths = %#v, want nil", root, paths)
	}
}
