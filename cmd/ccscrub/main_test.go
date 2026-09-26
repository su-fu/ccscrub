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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := run(tt.args); got != tt.want {
				t.Errorf("run(%q) = %d, want %d", tt.args, got, tt.want)
			}
		})
	}
}
