package protocol

import (
	"cod-server/internal/utils"
	"reflect"
	"testing"
	"time"
)

func TestEventCreation(t *testing.T) {
	method := "test_method"
	timestamp := time.Now().UTC()
	payload := utils.Dict{"key": "value", "number": 123}

	event := Event{
		Method:		method,
		Timestamp:	timestamp,
		Payload:	payload,
	}

	if event.Method != method {
		t.Errorf("Expected method %s, got %s", method, event.Method)
	}
	if !event.Timestamp.Equal(timestamp) {
		t.Errorf("Expected timestamp %v, got %v", timestamp, event.Timestamp)
	}
	if !reflect.DeepEqual(event.Payload, payload) {
		t.Errorf("Expected payload %v, got %v", payload, event.Payload)
	}
}
