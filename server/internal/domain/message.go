package domain

import (
	"time"
)

type MessageType string 

const (
	MessageTypeChat    MessageType = "chat"
	MessageTypeCommand MessageType = "command"
)

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	SenderID  string    `json:"sender_id"`
	RoomID    string 	  `json:"room_id"`
	Content   string    `json:"content"`
	Type      MessageType `json:"type"`
}
