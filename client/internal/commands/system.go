package commands

import (
	"cod-client/internal/api/protocol"
	"cod-client/internal/ui"
	"os"
	"strings"
)

func CommandClear(chat *ui.Chat, args []string)*protocol.Message {
	chat.Clear()
	return &protocol.Message{}
}

func CommandHelp(chat *ui.Chat, args []string)*protocol.Message {
	helpText := `Available commands:
/clear - Clear the chat window
/help - Show this help message
/exit - Exit the application
`
	chat.Write(helpText)	
	return &protocol.Message{}
}

func CommandExit(chat *ui.Chat, args []string)*protocol.Message {
	chat.Write("Exiting chat...")
	os.Exit(0)
	return &protocol.Message{}
}

func CommandUnknown(chat *ui.Chat, args []string)*protocol.Message {
	chat.Write("Unknown command: " + strings.Join(args, " "))
	return &protocol.Message{}
}
