package api

import (
	// "cod-server/internal/utils"
	"encoding/json"
	"time"
)

type Event struct {
	Method    string         `json:"method"`
	Timestamp time.Time      `json:"timestamp"`
	Payload   map[string]any `json:"payload"`
}

func (event Event) Json() ([]byte, error) {
	return json.Marshal(event)
}

func FromJson(data []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
