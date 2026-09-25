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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := run(tt.args); got != tt.want {
				t.Errorf("run(%q) = %d, want %d", tt.args, got, tt.want)
			}
		})
	}
}
