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
