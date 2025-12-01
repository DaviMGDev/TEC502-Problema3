package domain 

type MessageType uint8 

const (
	ChatMessage MessageType = iota 
	ServerNotification 
	GameEvent
	SystemAlert
)

type Message struct {
	ID      string      `json:"id"`
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
	Sender  string      `json:"sender"`
}
