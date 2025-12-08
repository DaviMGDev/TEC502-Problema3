package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	// "github.com/google/uuid"
)

type GameService interface {
	StartGame(player1, player2 string) (string, error)
	MakeMove(gameID, player, cardID string) error
	GetGameState(gameID string) (domain.Match, error)
	PLayerSurrender(gameID, player string) error
}

type GameServiceImplementation struct {
	gameRepo data.Repository[*domain.Match]
	userRepo data.Repository[*domain.User]
	matchCounter uint64
	queue chan *domain.Match
}

func NewGameService(gameRepo data.Repository[*domain.Match], userRepo data.Repository[*domain.User]) *GameServiceImplementation {
	return &GameServiceImplementation{
		gameRepo: gameRepo,
		userRepo: userRepo,
	}
}

func (s *GameServiceImplementation) StartGame(player1, player2 string) (string, error) {
	// matchID := uuid.New().String()
	// match := &domain.Match{}
	// match.ID = matchID 
	// match.Players = [2]string{player1, player2}
	// match.Winner = ""
	// match.Moves = [2]*domain.Card{}
	// err := s.gameRepo.Create(matchID, match)
	//
	//
	// if err != nil {
	// 	return "", err
	// }
	return "", nil
}

func (s *GameServiceImplementation) MakeMove(gameID, player, cardID string) error {
	// Implementation here
	return nil
}

func (s *GameServiceImplementation) GetGameState(gameID string) (domain.Match, error) {
	// Implementation here
	var match domain.Match
	return match, nil
}

func (s *GameServiceImplementation) PLayerSurrender(gameID, player string) error {
	// Implementation here
	return nil
}
