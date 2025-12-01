package protocol

import (
	"time"
)

type MessageType uint8

const (
)

type Message struct {
	Type MessageType `json:"type"`
	Timestamp time.Time   `json:"timestamp"`

}
