package main

import "testing"

func TestRun(t *testing.T) {
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
