package services

import (
	"cod-server/internal/domain"
	"cod-server/internal/data"
)

type CardServiceInterface interface {
	AddCard(userID string, card domain.Card) (error)
	GetCards(userID string) ([]domain.Card, error)
	RemoveCard(userID, cardID string) (error)
	ExchangeCards(userOneID, cardOneID, userTwoID, cardTwoID string) (error)
	LevelUpCard(userID, cardID string) (error)
}

type CardService struct {
	userRepo data.Repository[domain.User]
}

