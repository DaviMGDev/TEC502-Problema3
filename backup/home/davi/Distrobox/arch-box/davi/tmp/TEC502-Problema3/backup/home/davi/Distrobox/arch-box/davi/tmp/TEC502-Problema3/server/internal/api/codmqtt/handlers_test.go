package codmqtt

import (
	"cod-server/internal/api/protocol"
	"cod-server/internal/utils"
	"encoding/json"
	"testing"
	"time" // Added this import

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MockMQTTClient is a mock implementation of mqtt.Client for testing
type MockMQTTClient struct {
	PublishCalled bool
	PublishTopic  string
	PublishPayload []byte
	ConnectCalled bool
	SubscribeCalled bool
	SubscribeMultipleCalled bool
	// Add other necessary mock fields for client methods if needed
}

func (m *MockMQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	m.PublishCalled = true
	m.PublishTopic = topic
	if p, ok := payload.([]byte); ok {
		m.PublishPayload = p
	} else if s, ok := payload.(string); ok {
		m.PublishPayload = []byte(s)
	}
	return &MockToken{}
}
func (m *MockMQTTClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	m.SubscribeCalled = true
	return &MockToken{}
}
func (m *MockMQTTClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	m.SubscribeMultipleCalled = true
	return &MockToken{}
}
func (m *MockMQTTClient) Unsubscribe(topics ...string) mqtt.Token { return &MockToken{} }
func (m *MockMQTTClient) Disconnect(quiesce uint) {}
func (m *MockMQTTClient) IsConnected() bool { return true }
func (m *MockMQTTClient) IsConnectionOpen() bool { return true }

func (m *MockMQTTClient) OptionsReader() mqtt.ClientOptionsReader { return mqtt.ClientOptionsReader{} } 
func (m *MockMQTTClient) AddRoute(topic string, callback mqtt.MessageHandler) {}
func (m *MockMQTTClient) Connect() mqtt.Token {
	m.ConnectCalled = true
	return &MockToken{}
}

// MockToken is a mock implementation of mqtt.Token for testing
type MockToken struct {
	done chan struct{}
}

func (m *MockToken) Wait() bool { return true }
func (m *MockToken) WaitTimeout(duration time.Duration) bool { return true } // Corrected signature
func (m *MockToken) Error() error { return nil }
func (m *MockToken) Done() <-chan struct{} {
	if m.done == nil {
		m.done = make(chan struct{}, 1) // Buffered channel to prevent deadlock
		m.done <- struct{}{}
		close(m.done)
	}
	return m.done
}

// MockMessage is a mock implementation of mqtt.Message for testing
type MockMessage struct {
	payload []byte
	topic   string
}

func (m *MockMessage) Duplicate() bool { return false }
func (m *MockMessage) Qos() byte     { return 0 }
func (m *MockMessage) Retained() bool { return false }
func (m *MockMessage) Topic() string { return m.topic }
func (m *MockMessage) MessageID() uint16 { return 0 }
func (m *MockMessage) Payload() []byte { return m.payload }
func (m *MockMessage) Ack() {}

// MockHandlers for testing the wrap function
type MockBusinessLogicHandlers struct {
	OnRegisterEventCalled bool
}

func (m *MockBusinessLogicHandlers) OnRegisterEvent(event protocol.Event) protocol.Event {
	m.OnRegisterEventCalled = true
	return protocol.Event{
		Method:  "register",
		Payload: utils.Dict{"status": "success", "user_id": "mockUser"},
	}
}
func (m *MockBusinessLogicHandlers) OnLoginEvent(event protocol.Event) protocol.Event { return protocol.Event{} }
func (m *MockBusinessLogicHandlers) OnBuyCardPackEvent(event protocol.Event) protocol.Event { return protocol.Event{} }
func (m *MockBusinessLogicHandlers) OnTradeEvent(event protocol.Event) protocol.Event { return protocol.Event{} } // Changed from OnSwapCardEvent to OnTradeEvent
func (m *MockBusinessLogicHandlers) OnListUserCardsEvent(event protocol.Event) protocol.Event { return protocol.Event{} }
func (m *MockBusinessLogicHandlers) OnStartMatchEvent(event protocol.Event) protocol.Event { return protocol.Event{} }
func (m *MockBusinessLogicHandlers) OnPlayCardEvent(event protocol.Event) protocol.Event { return protocol.Event{} }
func (m *MockBusinessLogicHandlers) OnGetMatchResultEvent(event protocol.Event) protocol.Event { return protocol.Event{} }
func (m *MockBusinessLogicHandlers) OnPlayerSurrenderEvent(event protocol.Event) protocol.Event { return protocol.Event{} } // Changed from OnPLayerSurrenderEvent to OnPlayerSurrenderEvent


func TestWrapCallsBusinessLogicAndPublishesResponse(t *testing.T) {
	mockClient := &MockMQTTClient{}
	mockBusinessLogic := &MockBusinessLogicHandlers{}
	mh := NewMQTTHandler(mockClient, mockBusinessLogic)

	// Prepare a mock incoming MQTT message
	reqPayload := protocol.Event{
		Method:  "register",
		Payload: utils.Dict{"username": "test", "password": "pass", "client_id": "testclient"},
	}
	reqBytes, _ := json.Marshal(reqPayload)
	mockMsg := &MockMessage{payload: reqBytes, topic: "cod/request/register"}

	// Get the wrapped handler for OnRegisterEvent
	wrappedHandler := mh.wrap(mockBusinessLogic.OnRegisterEvent)

	// Simulate receiving an MQTT message
	wrappedHandler(mockClient, mockMsg)

	// Assertions
	if !mockBusinessLogic.OnRegisterEventCalled {
		t.Error("Expected OnRegisterEvent to be called on business logic handler, but it was not.")
	}
	if !mockClient.PublishCalled {
		t.Error("Expected MQTT client Publish to be called, but it was not.")
	}
	if mockClient.PublishTopic == "" {
		t.Errorf("Expected publish topic to be determined, got empty")
	}
	
	var resEvent protocol.Event
	json.Unmarshal(mockClient.PublishPayload, &resEvent)

	if resEvent.Method != "register" {
		t.Errorf("Expected response method 'register', got '%s'", resEvent.Method)
	}
	if status, ok := resEvent.Payload["status"].(string); !ok || status != "success" {
		t.Errorf("Expected response status 'success', got '%v'", resEvent.Payload["status"])
	}
	if userID, ok := resEvent.Payload["user_id"].(string); !ok || userID != "mockUser" {
		t.Errorf("Expected user_id 'mockUser', got '%v'", resEvent.Payload["user_id"])
	}

	expectedTopic := InferEventTopic(resEvent)
	if mockClient.PublishTopic != expectedTopic {
		t.Errorf("Expected publish topic '%s', got '%s'", expectedTopic, mockClient.PublishTopic)
	}
}

func TestInferEventTopic(t *testing.T) {
	tests := []struct {
		name          string
		event         protocol.Event
		expectedTopic string
	}{
		{
			name: "Register response topic",
			event: protocol.Event{
				Method:  "register",
				Payload: utils.Dict{"client_id": "client123"},
			},
			expectedTopic: "cod/response/users/register/client123",
		},
		{
			name: "Login response topic",
			event: protocol.Event{
				Method:  "login",
				Payload: utils.Dict{"client_id": "clientXYZ"},
			},
			expectedTopic: "cod/response/users/login/clientXYZ",
		},
		{
			name: "Start game response topic",
			event: protocol.Event{
				Method:  "start_game",
				Payload: utils.Dict{"user_id": "user456"},
			},
			expectedTopic: "cod/response/game/start_game/user456",
		},
		{
			name: "Play response topic",
			event: protocol.Event{
				Method:  "play",
				Payload: utils.Dict{"match_id": "matchABC", "user_id": "user789"},
			},
			expectedTopic: "cod/response/game/play/matchABC/user789",
		},
		{
			name: "Unknown method topic",
			event: protocol.Event{
				Method: "unknown_method",
			},
			expectedTopic: "cod/response/unknown_feature/unknown_method/unknown_identifier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topic := InferEventTopic(tt.event)
			if topic != tt.expectedTopic {
				t.Errorf("Expected topic '%s', got '%s'", tt.expectedTopic, topic)
			}
		})
	}
}
