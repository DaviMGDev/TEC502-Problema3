package api

import (
	"cod-client/internal/ui"
	"cod-client/internal/utils"
)

func HandleCommand(chat *ui.Chat, command string, args []string) {
	if handler, exists := table[command]; exists {
		handler(chat, args)
	} else {
		chat.Write("Unknown command: " + command)
	}
}
