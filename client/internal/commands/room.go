package commands

import (
	"cod-client/internal/api/protocol"
	"cod-client/internal/state"
	"cod-client/internal/ui"
	"cod-client/internal/utils"
	"time"
)

func CommandJoin(chat *ui.Chat, args []string)*protocol.Message {
	if len(args) < 1 {
	chat.Write("Usage: /join <room_id>")
		return &protocol.Message{}
	}
	roomID := args[0]
	return &protocol.Message{
		Method: "join_room",
		Timestamp: time.Now(),
		Status: "",
		Data: utils.Dict{
			"user_id": state.UserID,
			"room_id": roomID,
		},
	}
}

func CommandLeave(chat *ui.Chat, args []string)*protocol.Message {
	return &protocol.Message{
		Method: "leave_room",
		Timestamp: time.Now(),
		Status: "",
		Data: utils.Dict{
			"user_id": state.UserID,
			"room_id": state.RoomID,
		},
	}
}


