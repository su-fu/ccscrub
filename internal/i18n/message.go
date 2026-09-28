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
)

var messageEn = map[MessageID]string{
	MsgUsage:             "usage: ccscrub <scan|fix> [flags]",
	MsgUnknownSubcommand: "unknown subcommand: %q",
	MsgPathNotFound:      "%q not found",
	MsgPathAccessError:   "cannot access %q: %v",
}

var messageJa = map[MessageID]string{
	MsgUsage:             "使い方: ccscrub <scan|fix> [flags]",
	MsgUnknownSubcommand: "不明なサブコマンド: %q",
	MsgPathNotFound:      "%q が見つかりません",
	MsgPathAccessError:   "%q にアクセスできません: %v",
}

var messageZh = map[MessageID]string{
	MsgUsage:             "用法: ccscrub <scan|fix> [flags]",
	MsgUnknownSubcommand: "未知子命令: %q",
	MsgPathNotFound:      "未找到 %q",
	MsgPathAccessError:   "无法访问 %q: %v",
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
