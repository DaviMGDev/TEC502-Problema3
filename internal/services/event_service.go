package services

import (
	"cod-client/internal/api/protocol"
	"cod-client/internal/state"
	"cod-client/internal/utils"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// EventService encapsula a lógica de criação e publicação de eventos.
type EventService struct {
	appState *state.State
}

// NewEventService cria uma nova instância de EventService.
func NewEventService(s *state.State) *EventService {
	return &EventService{appState: s}
}

// createEvent é um helper genérico para criar um novo evento.
func (s *EventService) createEvent(method string, payload utils.Dict) protocol.Event {
	return protocol.Event{
		Method:    method,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

// inferTopicForEvent determina o tópico MQTT para um determinado evento.
func (s *EventService) inferTopicFor(event protocol.Event) string {
	switch event.Method {
	case "register", "login":
		return "user/" + event.Method + "/requests"
	case "chat":
		return "chat/room/" + s.appState.RoomID
	default:
		return ""
	}
}

// Publish serializa e publica um evento no tópico apropriado.
func (s *EventService) Publish(event protocol.Event) error {
	topic := s.inferTopicFor(event)
	if topic == "" {
		return fmt.Errorf("unknown topic for method: %s", event.Method)
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	token := s.appState.Client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

// --- Funções específicas de criação de eventos ---

func (s *EventService) CreateChatEvent(args []string) protocol.Event {
	return s.createEvent("chat", utils.Dict{
		"content": strings.Join(args, " "),
		"user_id": s.appState.UserID,
	})
}

func (s *EventService) CreateRegisterEvent(args []string) protocol.Event {
	return s.createEvent("register", utils.Dict{
		"username": args[0],
		"password": args[1],
	})
}

func (s *EventService) CreateLoginEvent(args []string) protocol.Event {
	return s.createEvent("login", utils.Dict{
		"username": args[0],
		"password": args[1],
	})
}