package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"errors"
	"strconv"
)

type GameService interface {
	StartGame(playerID string) (*domain.Match, error)
	MakeMove(gameID, playerID, cardID string) error
	GetGameState(gameID string) (*domain.Match, error)
	PlayerSurrender(gameID, player string) error
}

type GameServiceImplementation struct {
	gameRepo	data.Repository[*domain.Match]
	userRepo	data.Repository[*domain.User]
	matchCounter	uint64

	queue	chan *domain.Match
}

func NewGameService(gameRepo data.Repository[*domain.Match], userRepo data.Repository[*domain.User]) *GameServiceImplementation {
	return &GameServiceImplementation{
		gameRepo:	gameRepo,
		userRepo:	userRepo,

		queue:	make(chan *domain.Match, 1),
	}
}

func (s *GameServiceImplementation) StartGame(playerID string) (*domain.Match, error) {

	_, err := s.userRepo.Read(playerID)
	if err != nil {
		return nil, errors.New("player not found")
	}

	select {
	case match := <-s.queue:

		match.Players[1] = playerID
		err := s.gameRepo.Update(match.ID, match)
		if err != nil {

			return nil, err
		}
		return match, nil
	default:

		matchID := strconv.FormatUint(s.matchCounter, 10)
		s.matchCounter++

		match := &domain.Match{
			ID:		matchID,
			Players:	[2]string{playerID, ""},
		}

		err := s.gameRepo.Create(matchID, match)
		if err != nil {
			return nil, err
		}

		s.queue <- match

		return nil, nil
	}
}

func (s *GameServiceImplementation) MakeMove(gameID, playerID, cardID string) error {
	match, err := s.gameRepo.Read(gameID)
	if err != nil {
		return err
	}

	if match.Winner != "" {
		return errors.New("match has already concluded")
	}

	player, err := s.userRepo.Read(playerID)
	if err != nil {
		return errors.New("player not found")
	}

	card, ok := player.Cards.Get(cardID)
	if !ok {

		return errors.New("player does not own a card of this type")
	}

	playerIndex := -1
	if match.Players[0] == playerID {
		playerIndex = 0
	} else if match.Players[1] == playerID {
		playerIndex = 1
	} else {
		return errors.New("player is not part of this match")
	}

	if match.Moves[playerIndex] != nil {
		return errors.New("player has already made a move")
	}
	match.Moves[playerIndex] = card

	opponentIndex := (playerIndex + 1) % 2
	if match.Moves[opponentIndex] != nil {
		playerCard := match.Moves[playerIndex]
		opponentCard := match.Moves[opponentIndex]

		result := playerCard.Against(opponentCard)
		if result == 1 {
			match.Winner = match.Players[playerIndex]
		} else if result == -1 {
			match.Winner = match.Players[opponentIndex]
		} else {
			match.Winner = "draw"
		}
	}

	return s.gameRepo.Update(gameID, match)
}

func (s *GameServiceImplementation) GetGameState(gameID string) (*domain.Match, error) {
	return s.gameRepo.Read(gameID)
}

func (s *GameServiceImplementation) PlayerSurrender(gameID, player string) error {
	match, err := s.gameRepo.Read(gameID)
	if err != nil {
		return err
	}

	if match.Players[0] == player {
		match.Winner = match.Players[1]
	} else if match.Players[1] == player {
		match.Winner = match.Players[0]
	} else {
		return errors.New("player not in this match")
	}

	return s.gameRepo.Update(gameID, match)
}
