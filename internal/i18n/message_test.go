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
		{
			name: "skipped active(en)",
			lang: "en",
			id:   MsgSkippedActive,
			args: []any{"project", 4 * time.Minute},
			want: `[SKIP] "project" may be in use (modified 4m0s ago)`,
		},
		{
			name: "skipped active(ja)",
			lang: "ja",
			id:   MsgSkippedActive,
			args: []any{"project", 4 * time.Minute},
			want: `[SKIP] "project" は使用中の可能性があります (4m0s 前に更新)`,
		},
		{
			name: "skipped active(zh)",
			lang: "zh",
			id:   MsgSkippedActive,
			args: []any{"project", 4 * time.Minute},
			want: `[SKIP] "project" 可能正在使用 (4m0s 前更新)`,
		},
		{
			name: "skipped active(fr)",
			lang: "fr",
			id:   MsgSkippedActive,
			args: []any{"project", 4 * time.Minute},
			want: `[SKIP] "project" may be in use (modified 4m0s ago)`,
		},
		{
			name: "summary(en)",
			lang: "en",
			id:   MsgSummary,
			args: []any{2, 1},
			want: `2 secrets found in 1 files`,
		},
		{
			name: "summary(ja)",
			lang: "ja",
			id:   MsgSummary,
			args: []any{2, 1},
			want: `2 件のシークレットが 1 個のファイルで見つかりました`,
		},
		{
			name: "summary(zh)",
			lang: "zh",
			id:   MsgSummary,
			args: []any{2, 1},
			want: `发现 2 处机密信息，涉及 1 个文件`,
		},
		{
			name: "summary(fr)",
			lang: "fr",
			id:   MsgSummary,
			args: []any{2, 1},
			want: `2 secrets found in 1 files`,
		},
		{
			name: "fix summary(en)",
			lang: "en",
			id:   MsgFixSummary,
			args: []any{2, 1},
			want: `2 secrets redacted in 1 files`,
		},
		{
			name: "fix summary(ja)",
			lang: "ja",
			id:   MsgFixSummary,
			args: []any{2, 1},
			want: `2 件のシークレットを 1 個のファイルで更新しました`,
		},
		{
			name: "fix summary(zh)",
			lang: "zh",
			id:   MsgFixSummary,
			args: []any{2, 1},
			want: `已脱敏 2 处机密信息，涉及 1 个文件`,
		},
		{
			name: "fix summary(fr)",
			lang: "fr",
			id:   MsgFixSummary,
			args: []any{2, 1},
			want: `2 secrets redacted in 1 files`,
		},
		{
			name: "min-age zero(en)",
			lang: "en",
			id:   MsgMinAgeZero,
			args: nil,
			want: `warning: -min-age=0 also rewrites sessions that may be in use`,
		},
		{
			name: "min-age zero(ja)",
			lang: "ja",
			id:   MsgMinAgeZero,
			args: nil,
			want: `警告: -min-age=0 では、使用中の可能性があるセッションも書き換えます`,
		},
		{
			name: "min-age zero(zh)",
			lang: "zh",
			id:   MsgMinAgeZero,
			args: nil,
			want: `警告: -min-age=0 会同时改写可能正在使用的会话`,
		},
		{
			name: "min-age zero(fr)",
			lang: "fr",
			id:   MsgMinAgeZero,
			args: nil,
			want: `warning: -min-age=0 also rewrites sessions that may be in use`,
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
