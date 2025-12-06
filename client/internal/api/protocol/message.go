package protocol

import (
	"cod-client/internal/utils"
	"time"
)

type Message struct {
	Method string 		 	`json:"method"` 
	Status string      	`json:"status,omitempty"`
	Data utils.Dict  		`json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func NewMessage(method string, status string, data utils.Dict) *Message {
	return &Message{
		Method: method,
		Status: status,
		Data:   data,
		Timestamp: time.Now(),
	}
}

func (msg *Message) String() string {
	return utils.Dict{
		"method":    msg.Method,
		"status":    msg.Status,
		"data":      msg.Data,
		"timestamp": msg.Timestamp.Format(time.RFC3339),
	}.String()
}

func (msg *Message) Json() []byte {
	return utils.Dict{
		"method":    msg.Method,
		"status":    msg.Status,
		"data":      msg.Data,
		"timestamp": msg.Timestamp.Format(time.RFC3339),
	}.Json()
}

func ParseMessage(data []byte) (*Message, error) {
	var dict utils.Dict
	err := json.Unmarshal(data, &dict)
	if err != nil {
		return nil, err
	}
	method, _ := dict["method"].(string)
	status, _ := dict["status"].(string)
	dataField, _ := dict["data"].(utils.Dict)
	timestampStr, _ := dict["timestamp"].(string)
	timestamp, _ := time.Parse(time.RFC3339, timestampStr)

	return &Message{
		Method:    method,
		Status:    status,
		Data:      dataField,
		Timestamp: timestamp,
	}, nil
}	
