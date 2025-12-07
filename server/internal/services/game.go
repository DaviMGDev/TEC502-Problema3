package services

import (
	"cod-server/internal/domain"
	"cod-server/internal/data"
)

type GameServiceInterface interface {
	StartGame(roomID string) (error)
	PlayTurn(roomID, userID, cardID string) (error) 
	CalculateRoundWinner(roomID string) (error) 
	CalculateMatchWinner(roomID string) (error)
	GetMatchResult(roomID string) (domain.Match, error)
}

type GameService struct {
	matchRepo data.Repository[domain.Match]
	cardRepo  data.Repository[domain.Card]
	userRepo  data.Repository[domain.User]
}
