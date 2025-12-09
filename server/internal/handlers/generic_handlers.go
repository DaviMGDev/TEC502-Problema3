package handlers

import (
	"cod-server/internal/api/protocol"
	"cod-server/internal/services"
	"cod-server/internal/utils"
	"errors"
	"log"
)

type Handlers interface {
	OnRegisterEvent(event protocol.Event) protocol.Event
	OnLoginEvent(event protocol.Event) protocol.Event

	OnBuyCardPackEvent(event protocol.Event) protocol.Event
	OnSwapCardEvent(event protocol.Event) protocol.Event
	OnListUserCardsEvent(event protocol.Event) protocol.Event

	OnStartMatchEvent(event protocol.Event) protocol.Event
	OnPlayCardEvent(event protocol.Event) protocol.Event
	OnGetMatchResultEvent(event protocol.Event) protocol.Event
}

type HandlersImplementation struct {
	userService services.UserService
	cardService services.CardService
	gameService services.GameService
}

func NewHandlers(
	userService services.UserService,
	cardService services.CardService,
	gameService services.GameService,
) *HandlersImplementation {
	return &HandlersImplementation{
		userService: userService,
		cardService: cardService,
		gameService: gameService,
	}
}

// TODO: Implement these handlers with actual logic.
// For now, they just log the event and return a default response.

func (h *HandlersImplementation) OnRegisterEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnRegisterEvent: %+v", event)

	username, ok := event.Payload["username"].(string)
	if !ok || username == "" {
		return protocol.Event{
			Method: event.Method,
			Payload: utils.Dict{
				"status":         "error",
				"status_message": "missing or invalid username",
			},
		}
	}

	password, ok := event.Payload["password"].(string)
	if !ok || password == "" {
		return protocol.Event{
			Method: event.Method,
			Payload: utils.Dict{
				"status":         "error",
				"status_message": "missing or invalid password",
			},
		}
	}

	user, err := h.userService.Register(username, password)
	if err != nil {
		return protocol.Event{
			Method: event.Method,
			Payload: utils.Dict{
				"status":         "error",
				"status_message": err.Error(),
			},
		}
	}

	return protocol.Event{
		Method: event.Method,
		Payload: utils.Dict{
			"status":         "success",
			"status_message": "user registered successfully",
			"user_id":        user.ID,
		},
	}
}

func (h *HandlersImplementation) OnLoginEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnLoginEvent: %+v", event)
	return protocol.Event{}
}

func (h *HandlersImplementation) OnBuyCardPackEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnBuyCardPackEvent: %+v", event)
	return protocol.Event{}
}

func (h *HandlersImplementation) OnSwapCardEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnSwapCardEvent: %+v", event)
	return protocol.Event{}
}

func (h *HandlersImplementation) OnListUserCardsEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnListUserCardsEvent: %+v", event)
	return protocol.Event{}
}

func (h *HandlersImplementation) OnStartMatchEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnStartMatchEvent: %+v", event)

	playerID, ok := event.Payload["user_id"].(string)
	if !ok || playerID == "" {
		return protocol.Event{
			Method: event.Method,
			Payload: utils.Dict{
				"status":         "error",
				"status_message": "missing or invalid user_id",
			},
		}
	}

	match, err := h.gameService.StartGame(playerID)
	if err != nil {
		return protocol.Event{
			Method: event.Method,
			Payload: utils.Dict{
				"status":         "error",
				"status_message": err.Error(),
			},
		}
	}

	if match != nil {
		// Match formed
		return protocol.Event{
			Method: event.Method,
			Payload: utils.Dict{
				"status":         "success",
				"status_message": "match formed",
				"match_id":       match.ID,
			},
		}
	} else {
		// Player queued
		return protocol.Event{
			Method: event.Method,
			Payload: utils.Dict{
				"status":         "success",
				"status_message": "player queued",
			},
		}
	}
}

func (h *HandlersImplementation) OnPlayCardEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnPlayCardEvent: %+v", event)
	return protocol.Event{}
}

func (h *HandlersImplementation) OnGetMatchResultEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnGetMatchResultEvent: %+v", event)
	return protocol.Event{}
}
