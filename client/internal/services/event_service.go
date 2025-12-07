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
		return "user/" + event.Method
	case "chat":
		return "chat/room/" + s.appState.RoomID
	case "start":
		return "game/start_game"
	case "play":
		return "game/" + s.appState.RoomID + "/play_card"
	case "surrender":
		return "game/" + s.appState.RoomID + "/surrender"
	case "join":
		return "game/join_game"
	case "buy":
		return "store/buy"
	case "exchange":
		return "cards/" + s.appState.RoomID + "/exchange" + s.appState.UserID
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

func (s *EventService) CreateStartGameEvent() protocol.Event {
	return s.createEvent("start", utils.Dict{
		"user_id": s.appState.UserID,
//		"room_id": s.appState.RoomID,
	})
}

func (s *EventService) CreatePlayCardEvent(cardID string) protocol.Event {
	return s.createEvent("play", utils.Dict{
		"user_id": s.appState.UserID,
		"room_id": s.appState.RoomID,
		"card_id": cardID,
	})
}

func (s *EventService) CreateSurrenderEvent() protocol.Event {
	return s.createEvent("surrender", utils.Dict{
		"user_id": s.appState.UserID,
		"room_id": s.appState.RoomID,
	})
}

func (s *EventService) CreateJoinGameEvent(roomID string) protocol.Event {
	return s.createEvent("join", utils.Dict{
		"user_id": s.appState.UserID,
		"room_id": roomID,
	})
}

func (s *EventService) CreateBuyEvent(itemID string) protocol.Event {
	return s.createEvent("buy", utils.Dict{
		"user_id": s.appState.UserID,
		"item_id": itemID,
	})
}

func (s *EventService) CreateExchangeEvent(cardIDs []string) protocol.Event {
	return s.createEvent("exchange", utils.Dict{
		"user_id":  s.appState.UserID,
		"room_id":  s.appState.RoomID,
		"card_ids": cardIDs,
	})
}	
