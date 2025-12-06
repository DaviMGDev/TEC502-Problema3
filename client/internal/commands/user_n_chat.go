package commands

import (
	"cod-client/internal/api/protocol"
	"cod-client/internal/ui"
	"cod-client/internal/utils"
	"time"
)

func CommandChat(chat *ui.Chat, args []string)*protocol.Message {
	return &protocol.Message{
		Method: "chat",
		Timestamp: time.Now(),
		Status: "",
		Data: utils.Dict{
			"message": args,
		},
	}
}

func CommandRegister(chat *ui.Chat, args []string)*protocol.Message {
	if len(args) < 1 {
		chat.Write("Usage: /register <username>")
		return &protocol.Message{}
	}
	username := args[0]
	password := args[1]
	return &protocol.Message{
		Method: "register",
		Timestamp: time.Now(),
		Status: "",
		Data: utils.Dict{
			"username": username,
			"password": password,
		},
	}
}

func CommandLogin(chat *ui.Chat, args []string)*protocol.Message {
	if len(args) < 1 {
		chat.Write("Usage: /login <username>")
		return &protocol.Message{}
	}
	username := args[0]
	password := args[1]
	return &protocol.Message{
		Method: "login",
		Timestamp: time.Now(),
		Status: "",
		Data: utils.Dict{
			"username": username,
			"password": password,
		},
	}
}

func CommandLogout(chat *ui.Chat, args []string)*protocol.Message {
	return &protocol.Message{}
}
