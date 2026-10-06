package i18n

import (
	"fmt"
)

type MessageID int

const (
	MsgUsage MessageID = iota
	MsgUnknownSubcommand
	MsgPathNotFound
	MsgPathAccessError
	MsgSkipped
	MsgCannotRead
	MsgInvalidMinAge
)

var messageEn = map[MessageID]string{
	MsgUsage:             "usage: ccscrub <scan|fix> [flags]",
	MsgUnknownSubcommand: "unknown subcommand: %q",
	MsgPathNotFound:      "%q not found",
	MsgPathAccessError:   "cannot access %q: %v",
	MsgSkipped:           "[SKIP] %q (reason: %s)",
	MsgCannotRead:        "[SKIP] cannot read %q: %v",
	MsgInvalidMinAge:     "invalid value %q for flag -min-age: must not be negative",
}

var messageJa = map[MessageID]string{
	MsgUsage:             "使い方: ccscrub <scan|fix> [flags]",
	MsgUnknownSubcommand: "不明なサブコマンド: %q",
	MsgPathNotFound:      "%q が見つかりません",
	MsgPathAccessError:   "%q にアクセスできません: %v",
	MsgSkipped:           "[SKIP] %q (理由: %s)",
	MsgCannotRead:        "[SKIP] %q を読めません: %v",
	MsgInvalidMinAge:     "-min-age の値 %q は不正です: 負の値は指定できません",
}

var messageZh = map[MessageID]string{
	MsgUsage:             "用法: ccscrub <scan|fix> [flags]",
	MsgUnknownSubcommand: "未知子命令: %q",
	MsgPathNotFound:      "未找到 %q",
	MsgPathAccessError:   "无法访问 %q: %v",
	MsgSkipped:           "[SKIP] %q (原因: %s)",
	MsgCannotRead:        "[SKIP] 无法读取 %q: %v",
	MsgInvalidMinAge:     "-min-age 的值 %q 无效: 不能为负数",
}

var messages = map[string]map[MessageID]string{
	"en": messageEn,
	"ja": messageJa,
	"zh": messageZh,
}

func T(lang string, messageID MessageID, args ...any) string {
	v, ok := messages[lang][messageID]
	if ok {
		return fmt.Sprintf(v, args...)
	}
	return fmt.Sprintf(messages["en"][messageID], args...)
}
