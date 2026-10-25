package detect

import (
	"slices"
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	tok := "gh" + "p_" + strings.Repeat("A", 36)
	aws := "AK" + "IA" + strings.Repeat("A", 16)
	pat := "github" + "_pat_" + strings.Repeat("A", 22) + "_" + strings.Repeat("A", 59)
	ant := "sk-ant-" + "api03-" + strings.Repeat("A", 93) + "AA"
	antAdmin := "sk-ant-" + "admin01-" + strings.Repeat("A", 93) + "AA"
	openaiProj := "sk-" + "proj-" + strings.Repeat("A", 58) + "T3BlbkFJ" + strings.Repeat("A", 58)
	openaiSvcacct := "sk-" + "svcacct-" + strings.Repeat("A", 74) + "T3BlbkFJ" + strings.Repeat("A", 74)
	openaiAdmin := "sk-" + "admin-" + strings.Repeat("A", 58) + "T3BlbkFJ" + strings.Repeat("A", 74)
	openaiLegacy := "sk-" + strings.Repeat("A", 20) + "T3BlbkFJ" + strings.Repeat("A", 20)
	openaiNil := "sk-" + "proj-" + strings.Repeat("A", 58) + strings.Repeat("A", 58)

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
		{name: "aws key", line: aws, want: []Match{
			{RuleID: "aws-access-key-id", Start: 0, End: 20},
		}},
		{name: "aws key in env", line: "AWS_ACCESS_KEY_ID=" + aws, want: []Match{
			{RuleID: "aws-access-key-id", Start: 18, End: 38},
		}},
		{name: "aws temporary key", line: "AS" + "IA" + strings.Repeat("A", 16), want: []Match{
			{RuleID: "aws-access-key-id", Start: 0, End: 20},
		}},
		{name: "aws lowercase", line: "AK" + "IA" + strings.Repeat("a", 16), want: nil},
		{name: "aws after newline escape", line: "test" + `\n` + aws, want: []Match{
			{RuleID: "aws-access-key-id", Start: 6, End: 26},
		}},
		{name: "aws then github", line: aws + " " + tok, want: []Match{
			{RuleID: "aws-access-key-id", Start: 0, End: 20},
			{RuleID: "github-token", Start: 21, End: 61},
		}},
		{name: "github pat", line: pat, want: []Match{
			{RuleID: "github-token", Start: 0, End: 93},
		}},
		{name: "github pat too short", line: pat[:len(pat)-1], want: nil},
		{name: "anthropic key", line: ant, want: []Match{
			{RuleID: "anthropic-api-key", Start: 0, End: 108},
		}},
		{name: "anthropic key too short", line: ant[:len(ant)-1], want: nil},
		{name: "anthropic admin key", line: antAdmin, want: []Match{
			{RuleID: "anthropic-admin-api-key", Start: 0, End: 110},
		}},
		{name: "anthropic admin key too short", line: antAdmin[:len(antAdmin)-1], want: nil},
		{name: "openai proj key", line: openaiProj, want: []Match{
			{RuleID: "openai-api-key", Start: 0, End: 132},
		}},
		{name: "openai svcacct key", line: openaiSvcacct, want: []Match{
			{RuleID: "openai-api-key", Start: 0, End: 167},
		}},
		{name: "openai admin key", line: openaiAdmin, want: []Match{
			{RuleID: "openai-api-key", Start: 0, End: 149},
		}},
		{name: "openai legacy key", line: openaiLegacy, want: []Match{
			{RuleID: "openai-api-key", Start: 0, End: 51},
		}},
		{name: "openai nil key", line: openaiNil, want: nil},
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
