package detect

import (
	"slices"
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	tok := "gh" + "p_" + strings.Repeat("A", 36)
	tests := []struct {
		name string
		line string
		want []Match
	}{
		{name: "no secret", line: "hello", want: nil},
		{name: "token only", line: tok, want: []Match{
			{RuleID: "github-token", Start: 0, End: 40},
		}},
		{name: "token in text", line: "key=" + tok + "\n", want: []Match{
			{RuleID: "github-token", Start: 4, End: 44},
		}},
		{name: "two tokens", line: tok + " " + tok, want: []Match{
			{RuleID: "github-token", Start: 0, End: 40},
			{RuleID: "github-token", Start: 41, End: 81},
		}},
		{name: "wrong prefix", line: "gh" + "x_" + strings.Repeat("A", 36), want: nil},
		{name: "too short", line: "gh" + "p_" + strings.Repeat("A", 10), want: nil},
		{name: "after newline escape", line: `\n` + tok, want: []Match{
			{RuleID: "github-token", Start: 2, End: 42},
		}},
		{name: "after tab escape", line: `\t` + tok, want: []Match{
			{RuleID: "github-token", Start: 2, End: 42},
		}},
		{name: "after carriage return escape", line: `\r` + tok, want: []Match{
			{RuleID: "github-token", Start: 2, End: 42},
		}},
		{name: "inside a word", line: "x" + tok, want: nil},
		{name: "too long", line: tok + "A", want: nil},
		{name: "in JSON quotes", line: `"` + tok + `"`, want: []Match{
			{RuleID: "github-token", Start: 1, End: 41},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Find([]byte(tt.line))
			if !slices.Equal(got, tt.want) {
				t.Errorf("Find(%q) got = %+v, want = %+v", tt.line, got, tt.want)
			}
		})
	}
}

func TestPreview(t *testing.T) {
	tok := "gh" + "p_" + strings.Repeat("A", 32) + "WXYZ"
	tests := []struct {
		name string
		line string
		m    Match
		want string
	}{
		{name: "github token", line: tok, m: Match{RuleID: "github-token", Start: 0, End: 40}, want: "ghp_****WXYZ"},
		{name: "token in text", line: "key" + tok, m: Match{RuleID: "github-token", Start: 3, End: 43}, want: "ghp_****WXYZ"},
		{name: "exactly 12 bytes", line: strings.Repeat("B", 12), m: Match{RuleID: "github-token", Start: 0, End: 12}, want: "****"},
		{name: "13 bytes", line: "ABCD" + "xxxxx" + "WXYZ", m: Match{RuleID: "github-token", Start: 0, End: 13}, want: "ABCD****WXYZ"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Preview([]byte(tt.line), tt.m)
			if got != tt.want {
				t.Errorf("Preview(%q) got = %q, want = %q", tt.line, got, tt.want)
			}
		})
	}
}
