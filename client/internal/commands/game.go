package commands

import (
	"cod-client/internal/ui"
	"cod-client/internal/api/protocol"
	"cod-client/internal/utils"
	"cod-client/internal/state"
	"time"
)

func CommandStart(chat *ui.Chat, args []string)*protocol.Message {
	return &protocol.Message{
		Method:    "start",
		Timestamp: time.Now(),
		Status:    "",
		Data: utils.Dict{
			"user_id": state.UserID,
		},
	}
}

func CommandSurrender(chat *ui.Chat, args []string)*protocol.Message {
	return &protocol.Message{
		Method:    "surrender",
		Timestamp: time.Now(),
		Status:    "",
		Data: utils.Dict{
			"user_id": state.UserID,
			"room_id": state.RoomID,
		},
	}
}

func CommandListCards(chat *ui.Chat, args []string)*protocol.Message {
	return &protocol.Message{
		Method:    "list_cards",
		Timestamp: time.Now(),
		Status:    "",
		Data: utils.Dict{
			"user_id": state.UserID,
		},
	}
}

func CommandPlayCard(chat *ui.Chat, args []string)*protocol.Message {
	if len(args) < 1 {
		chat.Write("Usage: /play <card | card_id>")
		return &protocol.Message{}
	}
	card := args[0]
	return &protocol.Message{
		Method:    "play_card",
		Timestamp: time.Now(),
		Status:    "",
		Data: utils.Dict{
			"user_id": state.UserID,
			"room_id": state.RoomID,
			"card":    card,
		},
	}
}
