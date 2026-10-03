package discover

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
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

	got, skippedList, err := Find(root, 0, time.Now())
	if err != nil {
		t.Fatalf("Find(%q) error = %v", root, err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("Find(%q) = %v, want %v", root, got, want)
	}
	if skippedList != nil {
		t.Errorf("Find(%q) skipped = %#v, want nil", root, skippedList)
	}
}

func TestFindMissingRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "nope")
	paths, skippedList, err := Find(root, 30*time.Minute, time.Now())
	if err == nil {
		t.Fatalf("Find(%q) error = %v, want an error", root, err)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Find(%q) error = %v, want fs.ErrNotExist", root, err)
	}
	if paths != nil {
		t.Errorf("Find(%q) paths = %#v, want nil", root, paths)
	}
	if skippedList != nil {
		t.Errorf("Find(%q) skipped = %#v, want nil", root, skippedList)
	}
}

func TestFindSkipsActive(t *testing.T) {
	root := t.TempDir()

	paths := []string{
		"projA/old.jsonl",
		"projA/new.jsonl",
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

	now := time.Date(2026, time.October, 1, 12, 0, 0, 0, time.UTC)

	oldPath := filepath.Join(root, paths[0])
	newPath := filepath.Join(root, paths[1])

	oldMtime := now.Add(-2 * time.Hour)
	newMtime := now.Add(-5 * time.Minute)

	err := os.Chtimes(oldPath, oldMtime, oldMtime)
	if err != nil {
		t.Fatalf("Chtimes(%q): %v", oldPath, err)
	}
	err = os.Chtimes(newPath, newMtime, newMtime)
	if err != nil {
		t.Fatalf("Chtimes(%q): %v", newPath, err)
	}

	want := []string{oldPath}

	got, skippedList, err := Find(root, 30*time.Minute, now)
	if err != nil {
		t.Fatalf("Find(%q) error = %v", root, err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("Find(%q) = %v, want %v", root, got, want)
	}
	if len(skippedList) != 1 {
		t.Fatalf("Find(%q) skipped %d entries, want 1: %#v", root, len(skippedList), skippedList)
	}
	if skippedList[0].Path != newPath {
		t.Errorf("Skipped path = %q, want %q", skippedList[0].Path, newPath)
	}
	if skippedList[0].Reason != SkipActive {
		t.Errorf("Skipped reason = %q, want %q", skippedList[0].Reason, SkipActive)
	}
	if !skippedList[0].LastModified.Equal(newMtime) {
		t.Errorf("Last modified = %v, want %v", skippedList[0].LastModified, newMtime)
	}
}
