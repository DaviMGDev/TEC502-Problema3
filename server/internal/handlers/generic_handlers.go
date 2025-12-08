package handlers

import (
	"cod-server/internal/api/protocol"
	"cod-server/internal/services"
	// "cod-server/internal/domain"
	// "cod-server/internal/state"
	// "time"
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


