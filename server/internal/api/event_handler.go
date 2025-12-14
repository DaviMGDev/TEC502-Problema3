package api

import (
	"cod-server/internal/services"
	"time"
)

// EventHandler é a implementação da EventHandlerInterface
type EventHandler struct {
	userService   services.UserServiceInterface
	cardsService  services.CardsServiceInterface
	matchService  services.MatchServiceInterface
}

// NewEventHandler cria uma nova instância de EventHandler
func NewEventHandler(
	userService services.UserServiceInterface,
	cardsService services.CardsServiceInterface,
	matchService services.MatchServiceInterface,
) EventHandlerInterface {
	return &EventHandler{
		userService:   userService,
		cardsService:  cardsService,
		matchService:  matchService,
	}
}

func (eh *EventHandler) OnRegister(event Event) Event {
	username, ok1 := event.Payload["username"].(string)
	password, ok2 := event.Payload["password"].(string)
	if !ok1 || !ok2 {
		return makeErrorEvent("register_fail", "invalid payload")
	}

	err := eh.userService.Register(username, password)
	if err != nil {
		return makeErrorEvent("register_fail", err.Error())
	}

	return Event{
		Method:    "register_ok",
		Timestamp: time.Now(),
		Payload:   map[string]any{"username": username},
	}
}

func (eh *EventHandler) OnLogin(event Event) Event {
	username, ok1 := event.Payload["username"].(string)
	password, ok2 := event.Payload["password"].(string)
	if !ok1 || !ok2 {
		return makeErrorEvent("login_fail", "invalid payload")
	}

	user, err := eh.userService.Login(username, password)
	if err != nil {
		return makeErrorEvent("login_fail", err.Error())
	}
	if user == nil {
		return makeErrorEvent("login_fail", "invalid credentials")
	}

	return Event{
		Method:    "login_ok",
		Timestamp: time.Now(),
		Payload:   map[string]any{"user": user},
	}
}

func (eh *EventHandler) OnGetCards(event Event) Event {
	userID, ok := event.Payload["user_id"].(string)
	if !ok {
		return makeErrorEvent("get_cards_fail", "invalid payload")
	}

	cards, err := eh.cardsService.GetCards(userID)
	if err != nil {
		return makeErrorEvent("get_cards_fail", err.Error())
	}

	return Event{
		Method:    "get_cards_ok",
		Timestamp: time.Now(),
		Payload:   map[string]any{"cards": cards},
	}
}

func (eh *EventHandler) OnBuyPack(event Event) Event {
	userID, ok := event.Payload["user_id"].(string)
	if !ok {
		return makeErrorEvent("buy_pack_fail", "invalid payload")
	}

	err := eh.cardsService.BuyPack(userID)
	if err != nil {
		return makeErrorEvent("buy_pack_fail", err.Error())
	}

	return Event{
		Method:    "buy_pack_ok",
		Timestamp: time.Now(),
		Payload:   map[string]any{"user_id": userID},
	}
}

func (eh *EventHandler) OnOfferTrade(event Event) Event {
	fromUserID, ok1 := event.Payload["from_user_id"].(string)
	toUserID, ok2 := event.Payload["to_user_id"].(string)
	cardID, ok3 := event.Payload["card_id"].(string)
	if !ok1 || !ok2 || !ok3 {
		return makeErrorEvent("offer_trade_fail", "invalid payload")
	}

	err := eh.cardsService.OfferTrade(fromUserID, toUserID, cardID)
	if err != nil {
		return makeErrorEvent("offer_trade_fail", err.Error())
	}
	return Event{
		Method:    "offer_trade_ok",
		Timestamp: time.Now(),
		Payload:   event.Payload,
	}
}

func (eh *EventHandler) OnStartMatch(event Event) Event {
	userID, ok := event.Payload["user_id"].(string)
	if !ok {
		return makeErrorEvent("start_match_fail", "invalid payload")
	}
	match, err := eh.matchService.StartMatch(userID)
	if err != nil {
		return makeErrorEvent("start_match_fail", err.Error())
	}

	return Event{
		Method:    "start_match_ok",
		Timestamp: time.Now(),
		Payload:   map[string]any{"match": match},
	}
}

func (eh *EventHandler) OnJoinMatch(event Event) Event {
	userID, ok1 := event.Payload["user_id"].(string)
	matchID, ok2 := event.Payload["match_id"].(string)
	if !ok1 || !ok2 {
		return makeErrorEvent("join_match_fail", "invalid payload")
	}

	err := eh.matchService.JoinMatch(userID, matchID)
	if err != nil {
		return makeErrorEvent("join_match_fail", err.Error())
	}

	return Event{
		Method:    "join_match_ok",
		Timestamp: time.Now(),
		Payload:   event.Payload,
	}
}

func (eh *EventHandler) OnSurrenderMatch(event Event) Event {
	userID, ok1 := event.Payload["user_id"].(string)
	matchID, ok2 := event.Payload["match_id"].(string)
	if !ok1 || !ok2 {
		return makeErrorEvent("surrender_match_fail", "invalid payload")
	}
	err := eh.matchService.SurrenderMatch(userID, matchID)
	if err != nil {
		return makeErrorEvent("surrender_match_fail", err.Error())
	}
	return Event{
		Method:    "surrender_match_ok",
		Timestamp: time.Now(),
		Payload:   event.Payload,
	}
}

func (eh *EventHandler) OnMakeMove(event Event) Event {
	userID, ok1 := event.Payload["user_id"].(string)
	matchID, ok2 := event.Payload["match_id"].(string)
	cardID, ok3 := event.Payload["card_id"].(string)
	if !ok1 || !ok2 || !ok3 {
		return makeErrorEvent("make_move_fail", "invalid payload")
	}

	err := eh.matchService.MakeMove(userID, matchID, cardID)
	if err != nil {
		return makeErrorEvent("make_move_fail", err.Error())
	}
	return Event{
		Method:    "make_move_ok",
		Timestamp: time.Now(),
		Payload:   event.Payload,
	}
}

// makeErrorEvent é uma função helper para criar eventos de erro padronizados
func makeErrorEvent(method, message string) Event {
	return Event{
		Method:    method,
		Timestamp: time.Now(),
		Payload:   map[string]any{"error": message},
	}
}


