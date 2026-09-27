package i18n

import "testing"

func TestT(t *testing.T) {
	tests := []struct {
		name string
		lang string
		id   MessageID
		args []any
		want string
	}{
		{
			name: "usage(en)",
			lang: "en",
			id:   MsgUsage,
			args: nil,
			want: `usage: ccscrub <scan|fix> [flags]`,
		},
		{
			name: "usage(ja)",
			lang: "ja",
			id:   MsgUsage,
			args: nil,
			want: `使い方: ccscrub <scan|fix> [flags]`,
		},
		{
			name: "usage(zh)",
			lang: "zh",
			id:   MsgUsage,
			args: nil,
			want: `用法: ccscrub <scan|fix> [flags]`,
		},
		{
			name: "usage(fr)",
			lang: "fr",
			id:   MsgUsage,
			args: nil,
			want: `usage: ccscrub <scan|fix> [flags]`,
		},
		{
			name: "unknown subcommand(en)",
			lang: "en",
			id:   MsgUnknownSubcommand,
			args: []any{"foo"},
			want: `unknown subcommand: "foo"`,
		},
		{
			name: "unknown subcommand(ja)",
			lang: "ja",
			id:   MsgUnknownSubcommand,
			args: []any{"foo"},
			want: `不明なサブコマンド: "foo"`,
		},
		{
			name: "unknown subcommand(zh)",
			lang: "zh",
			id:   MsgUnknownSubcommand,
			args: []any{"foo"},
			want: `未知子命令: "foo"`,
		},
		{
			name: "unknown subcommand(fr)",
			lang: "fr",
			id:   MsgUnknownSubcommand,
			args: []any{"foo"},
			want: `unknown subcommand: "foo"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := T(tt.lang, tt.id, tt.args...); got != tt.want {
				t.Errorf("T(%q, %d, %v) = %q, want %q", tt.lang, tt.id, tt.args, got, tt.want)
			}
		})
	}
}
