package protocol

import (
	"time"
	"cod-server/internal/utils"
)

type Event struct {
	Method		string		`json:"method"`
	Timestamp	time.Time	`json:"timestamp"`
	Payload		utils.Dict	`json:"payload"`
}
