package i18n

import (
	"fmt"
)

type MessageID int

const (
	MsgUsage             MessageID = iota // 0
	MsgUnknownSubcommand                  // 1
)

var messageEn = map[MessageID]string{
	MsgUsage:             "usage: ccscrub <scan|fix> [flags]",
	MsgUnknownSubcommand: "unknown subcommand: %q",
}

var messageJa = map[MessageID]string{
	MsgUsage:             "usage: ccscrub <scan|fix> [flags]",
	MsgUnknownSubcommand: "unknown subcommand: %q",
}

var messageZh = map[MessageID]string{
	MsgUsage:             "usage: ccscrub <scan|fix> [flags]",
	MsgUnknownSubcommand: "unknown subcommand: %q",
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
