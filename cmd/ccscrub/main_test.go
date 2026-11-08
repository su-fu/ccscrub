package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	setupHome(t)
	tests := []struct {
		name string
		args []string
		want int
	}{
		{
			name: "no args",
			args: []string{"ccscrub"},
			want: 2,
		},
		{
			name: "scan",
			args: []string{"ccscrub", "scan"},
			want: 0,
		},
		{
			name: "fix",
			args: []string{"ccscrub", "fix"},
			want: 0,
		},
		{
			name: "foo",
			args: []string{"ccscrub", "foo"},
			want: 2,
		},
		{
			name: "scan --foo",
			args: []string{"ccscrub", "scan", "--foo"},
			want: 2,
		},
		{
			name: "scan -h",
			args: []string{"ccscrub", "scan", "-h"},
			want: 0,
		},
		{
			name: "fix -h",
			args: []string{"ccscrub", "fix", "-h"},
			want: 0,
		},
		{
			name: "fix --foo",
			args: []string{"ccscrub", "fix", "--foo"},
			want: 2,
		},
		{
			name: "scan --min-age=-5m",
			args: []string{"ccscrub", "scan", "--min-age=-5m"},
			want: 2,
		},
		{
			name: "fix --min-age=-5m",
			args: []string{"ccscrub", "fix", "--min-age=-5m"},
			want: 2,
		},
		{
			name: "scan --min-age=abc",
			args: []string{"ccscrub", "scan", "--min-age=abc"},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := run(tt.args); got != tt.want {
				t.Errorf("run(%q) = %d, want %d", tt.args, got, tt.want)
			}
		})
	}
}

func TestResolveLang(t *testing.T) {
	tests := []struct {
		name     string
		flagLang string
		envLang  string
		want     string
	}{
		{
			name:     "flag ja, no env",
			flagLang: "ja",
			envLang:  "",
			want:     "ja",
		},
		{
			name:     "no flag, env zh",
			flagLang: "",
			envLang:  "zh",
			want:     "zh",
		},
		{
			name:     "flag ja, env zh",
			flagLang: "ja",
			envLang:  "zh",
			want:     "ja",
		},
		{
			name:     "no flag, no env",
			flagLang: "",
			envLang:  "",
			want:     "en",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveLang(tt.flagLang, tt.envLang); got != tt.want {
				t.Errorf("resolveLang(%q, %q) = %q, want %q", tt.flagLang, tt.envLang, got, tt.want)
			}
		})
	}
}

func TestRunScanFindsSecret(t *testing.T) {
	dir := setupHome(t)
	path := filepath.Join(dir, "test.jsonl")

	tok := "gh" + "p_" + strings.Repeat("A", 36)
	want := `{"text":"` + tok + `"}` + "\n"

	err := os.WriteFile(path, []byte(want), 0o644)
	if err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
	args := []string{"ccscrub", "scan", "--min-age=0"}
	got := run(args)

	if got != 1 {
		t.Errorf("run(%q) = %d, want %d", args, got, 1)
	}
}

func TestRedactLine(t *testing.T) {
	tok := "gh" + "p_" + strings.Repeat("A", 36)
	tests := []struct {
		name     string
		line     string
		wantLine string
		wantOK   bool
	}{
		{name: "valid", line: `{"text":"` + tok + `"}` + "\n", wantLine: `{"text":"[REDACTED:github-token]"}` + "\n", wantOK: true},
		{name: "broken", line: `{"text":"` + tok + "\n", wantLine: `{"text":"[REDACTED:github-token]` + "\n", wantOK: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLine, gotOK := redactLine([]byte(tt.line))
			if !(slices.Equal(gotLine, []byte(tt.wantLine)) && gotOK == tt.wantOK) {
				t.Errorf("redactLine(%q) = gotLine %q, gotOK %v, wantLine %q, wantOK %v", tt.line, string(gotLine), gotOK, tt.wantLine, tt.wantOK)
			}
		})
	}
}

func setupHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "projects")
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		t.Fatalf("MkdirAll(%q): %v", dir, err)
	}
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)

	return dir
}

func TestFixFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.jsonl")

	tokA := "gh" + "p_" + strings.Repeat("A", 36)
	tokB := "gh" + "p_" + strings.Repeat("B", 36)
	tokC := "gh" + "p_" + strings.Repeat("C", 36)
	input := `{"text":"` + tokA + `"}` + "\n" + `{"text":"` + tokB + `"}` + "\n" + `{"text":"` + tokC + `"}` + "\n" + `{"text":"hello"}` + "\n"
	want := `{"text":"[REDACTED:github-token]"}` + "\n" + `{"text":"[REDACTED:github-token]"}` + "\n" + `{"text":"[REDACTED:github-token]"}` + "\n" + `{"text":"hello"}` + "\n"

	err := os.WriteFile(path, []byte(input), 0o644)
	if err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}

	n, err := fixFile(path)
	if err != nil {
		t.Fatalf("fixFile(%q): %v", path, err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", path, err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", dir, err)
	}

	if n != 3 {
		t.Errorf("fixFile(%q) n = %d, want 3", path, n)
	}

	if !slices.Equal(got, []byte(want)) {
		t.Errorf("fixFile(%q) got = %q, want = %q", path, string(got), string(want))
	}

	if len(entries) != 1 {
		t.Errorf("ReadDir(%q) = %d entries, want 1", dir, len(entries))
	}
}

func TestFixFileNoSecrets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.jsonl")

	input := `{"text":"hello"}` + "\n"
	want := `{"text":"hello"}` + "\n"

	err := os.WriteFile(path, []byte(input), 0o644)
	if err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}

	old := time.Date(2026, time.October, 1, 11, 0, 0, 0, time.UTC)
	err = os.Chtimes(path, old, old)
	if err != nil {
		t.Fatalf("Chtimes(%q): %v", path, err)
	}

	n, err := fixFile(path)
	if err != nil {
		t.Fatalf("fixFile(%q): %v", path, err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", path, err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", dir, err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat(%q): %v", path, err)
	}

	if n != 0 {
		t.Errorf("fixFile(%q) n = %d, want 0", path, n)
	}

	if !slices.Equal(got, []byte(want)) {
		t.Errorf("fixFile(%q) got = %q, want = %q", path, string(got), string(want))
	}

	if len(entries) != 1 {
		t.Errorf("ReadDir(%q) = %d entries, want 1", dir, len(entries))
	}

	if !old.Equal(fi.ModTime()) {
		t.Errorf("ModTime = %v, want = %v", fi.ModTime(), old)
	}
}
