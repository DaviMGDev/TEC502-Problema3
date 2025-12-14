package services 

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
)

type CardsService struct {
	cardsRepo data.Repository[domain.CardInterface]
}

func NewCardsService(cardsRepo data.Repository[domain.CardInterface]) CardsServiceInterface {
	return &CardsService{cardsRepo: cardsRepo}
}

func (cs *CardsService) GetCards(userID string) ([]domain.CardInterface, error) {
	cards, err := cs.cardsRepo.ListBy(func(c domain.CardInterface) bool {
		return c.GetOwnerID() == userID
	})
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (cs *CardsService) BuyPack(userID string) error {
	// Implementation for buying a pack of cards
	return nil
}

func (cs *CardsService) OfferTrade(fromUserID, toUserID, cardID string) error {
	// Implementation for offering a trade
	return nil
}

func (cs *CardsService) AcceptTrade(fromUserID, toUserID, cardID string) error {
	// Implementation for accepting a trade
	return nil
}	

