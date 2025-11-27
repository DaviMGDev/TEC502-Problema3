package domain

import (
	"time"
)

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	SenderID  string    `json:"sender_id"`
	RoomID    string 	  `json:"room_id"`
	Content   string    `json:"content"`
}
