package services

import (
	"cod-server/internal/domain"
)

type GameService interface {
	StartGame(roomID string) (error)
	PlayTurn(roomID, userID string, cardType domain.CardType) (error) 
	CalculateRoundWinner(roomID string) (error) 
	CalculateMatchWinner(roomID string) (error)
	GetMatchResult(roomID string) (domain.Match, error)
}

