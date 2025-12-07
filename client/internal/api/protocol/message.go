package protocol

import (
	"cod-client/internal/utils"
	"time"
)

type Event struct {
	Method    string     `json:"method"`
	Timestamp time.Time  `json:"timestamp"`
	Payload   utils.Dict `json:"payload"`
}