package i18n

import (
	"errors"
	"testing"
	"time"
)

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
		{
			name: "path not found(en)",
			lang: "en",
			id:   MsgPathNotFound,
			args: []any{"projects"},
			want: `"projects" not found`,
		},
		{
			name: "path not found(ja)",
			lang: "ja",
			id:   MsgPathNotFound,
			args: []any{"projects"},
			want: `"projects" が見つかりません`,
		},
		{
			name: "path not found(zh)",
			lang: "zh",
			id:   MsgPathNotFound,
			args: []any{"projects"},
			want: `未找到 "projects"`,
		},
		{
			name: "path not found(fr)",
			lang: "fr",
			id:   MsgPathNotFound,
			args: []any{"projects"},
			want: `"projects" not found`,
		},
		{
			name: "path access error(en)",
			lang: "en",
			id:   MsgPathAccessError,
			args: []any{"projects", errors.New("permission denied")},
			want: `cannot access "projects": permission denied`,
		},
		{
			name: "path access error(ja)",
			lang: "ja",
			id:   MsgPathAccessError,
			args: []any{"projects", errors.New("permission denied")},
			want: `"projects" にアクセスできません: permission denied`,
		},
		{
			name: "path access error(zh)",
			lang: "zh",
			id:   MsgPathAccessError,
			args: []any{"projects", errors.New("permission denied")},
			want: `无法访问 "projects": permission denied`,
		},
		{
			name: "path access error(fr)",
			lang: "fr",
			id:   MsgPathAccessError,
			args: []any{"projects", errors.New("permission denied")},
			want: `cannot access "projects": permission denied`,
		},
		{
			name: "path cannot read(en)",
			lang: "en",
			id:   MsgCannotRead,
			args: []any{"projects", errors.New("permission denied")},
			want: `[SKIP] cannot read "projects": permission denied`,
		},
		{
			name: "path cannot read(ja)",
			lang: "ja",
			id:   MsgCannotRead,
			args: []any{"projects", errors.New("permission denied")},
			want: `[SKIP] "projects" を読めません: permission denied`,
		},
		{
			name: "path cannot read(zh)",
			lang: "zh",
			id:   MsgCannotRead,
			args: []any{"projects", errors.New("permission denied")},
			want: `[SKIP] 无法读取 "projects": permission denied`,
		},
		{
			name: "path cannot read(fr)",
			lang: "fr",
			id:   MsgCannotRead,
			args: []any{"projects", errors.New("permission denied")},
			want: `[SKIP] cannot read "projects": permission denied`,
		},
		{
			name: "skipped(en)",
			lang: "en",
			id:   MsgSkipped,
			args: []any{"projects", "active"},
			want: `[SKIP] "projects" (reason: active)`,
		},
		{
			name: "skipped(ja)",
			lang: "ja",
			id:   MsgSkipped,
			args: []any{"projects", "active"},
			want: `[SKIP] "projects" (理由: active)`,
		},
		{
			name: "skipped(zh)",
			lang: "zh",
			id:   MsgSkipped,
			args: []any{"projects", "active"},
			want: `[SKIP] "projects" (原因: active)`,
		},
		{
			name: "skipped(fr)",
			lang: "fr",
			id:   MsgSkipped,
			args: []any{"projects", "active"},
			want: `[SKIP] "projects" (reason: active)`,
		},
		{
			name: "invalid min-age(en)",
			lang: "en",
			id:   MsgInvalidMinAge,
			args: []any{-5 * time.Minute},
			want: `invalid value "-5m0s" for flag -min-age: must not be negative`,
		},
		{
			name: "invalid min-age(ja)",
			lang: "ja",
			id:   MsgInvalidMinAge,
			args: []any{-5 * time.Minute},
			want: `-min-age の値 "-5m0s" は不正です: 負の値は指定できません`,
		},
		{
			name: "invalid min-age(zh)",
			lang: "zh",
			id:   MsgInvalidMinAge,
			args: []any{-5 * time.Minute},
			want: `-min-age 的值 "-5m0s" 无效: 不能为负数`,
		},
		{
			name: "invalid min-age(fr)",
			lang: "fr",
			id:   MsgInvalidMinAge,
			args: []any{-5 * time.Minute},
			want: `invalid value "-5m0s" for flag -min-age: must not be negative`,
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
