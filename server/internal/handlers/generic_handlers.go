package handlers

import (
	"cod-server/internal/api/protocol"
	"cod-server/internal/services"
	"cod-server/internal/utils"
	_ "errors"
	"log"
)

type Handlers interface {
	OnRegisterEvent(event protocol.Event) protocol.Event
	OnLoginEvent(event protocol.Event) protocol.Event

	OnBuyCardPackEvent(event protocol.Event) protocol.Event
	OnTradeEvent(event protocol.Event) protocol.Event
	OnListUserCardsEvent(event protocol.Event) protocol.Event

	OnStartMatchEvent(event protocol.Event) protocol.Event
	OnPlayCardEvent(event protocol.Event) protocol.Event
	OnGetMatchResultEvent(event protocol.Event) protocol.Event
	OnPlayerSurrenderEvent(event protocol.Event) protocol.Event
}

type HandlersImplementation struct {
	userService	services.UserService
	cardService	services.CardService
	gameService	services.GameService
}

func NewHandlers(
	userService services.UserService,
	cardService services.CardService,
	gameService services.GameService,
) *HandlersImplementation {
	return &HandlersImplementation{
		userService:	userService,
		cardService:	cardService,
		gameService:	gameService,
	}
}

func (h *HandlersImplementation) OnRegisterEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnRegisterEvent: %+v", event)

	username, ok := event.Payload["username"].(string)
	if !ok || username == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid username",
			},
		}
	}

	password, ok := event.Payload["password"].(string)
	if !ok || password == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid password",
			},
		}
	}

	user, err := h.userService.Register(username, password)
	if err != nil {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	err.Error(),
			},
		}
	}

	return protocol.Event{
		Method:	event.Method,
		Payload: utils.Dict{
			"status":		"success",
			"status_message":	"user registered successfully",
			"user_id":		user.ID,
		},
	}
}

func (h *HandlersImplementation) OnLoginEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnLoginEvent: %+v", event)

	username, ok := event.Payload["username"].(string)
	if !ok || username == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid username",
			},
		}
	}

	password, ok := event.Payload["password"].(string)
	if !ok || password == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid password",
			},
		}
	}

	userID, err := h.userService.Login(username, password)
	if err != nil {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	err.Error(),
			},
		}
	}

	return protocol.Event{
		Method:	event.Method,
		Payload: utils.Dict{
			"status":		"success",
			"status_message":	"user logged in successfully",
			"user_id":		userID,
		},
	}
}

func (h *HandlersImplementation) OnBuyCardPackEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnBuyCardPackEvent: %+v", event)

	userID, ok := event.Payload["user_id"].(string)
	if !ok || userID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid user_id",
			},
		}
	}

	err := h.cardService.BuyCardPack(userID)
	if err != nil {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	err.Error(),
			},
		}
	}

	return protocol.Event{
		Method:	event.Method,
		Payload: utils.Dict{
			"status":		"success",
			"status_message":	"card pack bought successfully",
		},
	}
}

func (h *HandlersImplementation) OnTradeEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnTradeEvent: %+v", event)

	user1ID, ok := event.Payload["user_id"].(string)
	if !ok || user1ID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid user_id for trading user",
			},
		}
	}

	targetUserID, ok := event.Payload["target_user_id"].(string)
	if !ok || targetUserID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid target_user_id",
			},
		}
	}

	cardType, ok := event.Payload["card_type"].(string)
	if !ok || cardType == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid card_type",
			},
		}
	}

	err := h.cardService.Trade(user1ID, targetUserID, cardType)
	if err != nil {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	err.Error(),
			},
		}
	}

	return protocol.Event{
		Method:	event.Method,
		Payload: utils.Dict{
			"status":		"success",
			"status_message":	"cards traded successfully",
		},
	}
}

func (h *HandlersImplementation) OnListUserCardsEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnListUserCardsEvent: %+v", event)

	userID, ok := event.Payload["user_id"].(string)
	if !ok || userID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid user_id",
			},
		}
	}

	cards, err := h.cardService.ListUserCards(userID)
	if err != nil {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	err.Error(),
			},
		}
	}

	cardList := make([]map[string]any, len(cards))
	for i, card := range cards {
		cardList[i] = map[string]any{
			"id":		card.ID,
			"type":		card.Type,
			"level":	card.Level,
		}
	}

	return protocol.Event{
		Method:	event.Method,
		Payload: utils.Dict{
			"status":		"success",
			"status_message":	"user cards listed successfully",
			"cards":		cardList,
		},
	}
}

func (h *HandlersImplementation) OnStartMatchEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnStartMatchEvent: %+v", event)

	playerID, ok := event.Payload["user_id"].(string)
	if !ok || playerID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid user_id",
			},
		}
	}

	match, err := h.gameService.StartGame(playerID)
	if err != nil {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	err.Error(),
			},
		}
	}

	if match != nil {

		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"success",
				"status_message":	"match formed",
				"match_id":		match.ID,
			},
		}
	} else {

		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"success",
				"status_message":	"player queued",
			},
		}
	}
}

func (h *HandlersImplementation) OnPlayCardEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnPlayCardEvent: %+v", event)

	gameID, ok := event.Payload["match_id"].(string)
	if !ok || gameID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid match_id",
			},
		}
	}

	playerID, ok := event.Payload["user_id"].(string)
	if !ok || playerID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid user_id",
			},
		}
	}

	cardID, ok := event.Payload["move"].(string)
	if !ok || cardID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid move (card type)",
			},
		}
	}

	err := h.gameService.MakeMove(gameID, playerID, cardID)
	if err != nil {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	err.Error(),
			},
		}
	}

	match, err := h.gameService.GetGameState(gameID)
	if err != nil {

		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"failed to retrieve game state after move: " + err.Error(),
			},
		}
	}

	responsePayload := utils.Dict{
		"status":		"success",
		"status_message":	"move made successfully",
	}

	if match.Winner != "" {
		responsePayload["round_result"] = match.Winner
	} else {
		responsePayload["status_message"] = "move made, waiting for opponent"
	}

	return protocol.Event{
		Method:		event.Method,
		Payload:	responsePayload,
	}
}

func (h *HandlersImplementation) OnPlayerSurrenderEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnPlayerSurrenderEvent: %+v", event)

	gameID, ok := event.Payload["match_id"].(string)
	if !ok || gameID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid match_id",
			},
		}
	}

	playerID, ok := event.Payload["user_id"].(string)
	if !ok || playerID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid user_id",
			},
		}
	}

	err := h.gameService.PlayerSurrender(gameID, playerID)
	if err != nil {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	err.Error(),
			},
		}
	}

	return protocol.Event{
		Method:	event.Method,
		Payload: utils.Dict{
			"status":		"success",
			"status_message":	"player surrendered successfully",
		},
	}
}

func (h *HandlersImplementation) OnGetMatchResultEvent(event protocol.Event) protocol.Event {
	log.Printf("Received OnGetMatchResultEvent: %+v", event)

	gameID, ok := event.Payload["match_id"].(string)
	if !ok || gameID == "" {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	"missing or invalid match_id",
			},
		}
	}

	match, err := h.gameService.GetGameState(gameID)
	if err != nil {
		return protocol.Event{
			Method:	event.Method,
			Payload: utils.Dict{
				"status":		"error",
				"status_message":	err.Error(),
			},
		}
	}

	moves := make([]map[string]any, 2)
	if match.Moves[0] != nil {
		moves[0] = map[string]any{"id": match.Moves[0].ID, "type": match.Moves[0].Type, "level": match.Moves[0].Level}
	} else {
		moves[0] = nil
	}
	if match.Moves[1] != nil {
		moves[1] = map[string]any{"id": match.Moves[1].ID, "type": match.Moves[1].Type, "level": match.Moves[1].Level}
	} else {
		moves[1] = nil
	}

	return protocol.Event{
		Method:	event.Method,
		Payload: utils.Dict{
			"status":		"success",
			"status_message":	"game state retrieved successfully",
			"match_id":		match.ID,
			"players":		match.Players,
			"moves":		moves,
			"winner":		match.Winner,
		},
	}
}
