package services

import (
	"cod-server/internal/domain"
)

type CardService interface {
	AddCard(userID string, card domain.Card) (error)
	GetCards(userID string) ([]domain.Card, error)
	RemoveCard(userID string, cardType domain.CardType) (error)
	ExchangeCards(userOneID, userTwoID string, cardType domain.CardType) (error)
	LevelUpCard(userID string, cardType domain.CardType) (error)
}


