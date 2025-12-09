package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"errors"
	"strconv"
)

type GameService interface {
	// StartGame returns a match if one is formed, otherwise returns nil.
	StartGame(playerID string) (*domain.Match, error)
	MakeMove(gameID, playerID, cardID string) error
	GetGameState(gameID string) (*domain.Match, error)
	PLayerSurrender(gameID, player string) error
}

type GameServiceImplementation struct {
	gameRepo     data.Repository[*domain.Match]
	userRepo     data.Repository[*domain.User]
	matchCounter uint64
	// A buffered channel of size 1 acts as our queue for one pending match
	queue chan *domain.Match
}

func NewGameService(gameRepo data.Repository[*domain.Match], userRepo data.Repository[*domain.User]) *GameServiceImplementation {
	return &GameServiceImplementation{
		gameRepo: gameRepo,
		userRepo: userRepo,
		// Buffer size of 1 allows one match to wait in the queue
		queue: make(chan *domain.Match, 1),
	}
}

func (s *GameServiceImplementation) StartGame(playerID string) (*domain.Match, error) {
	// Check if player exists
	_, err := s.userRepo.Read(playerID)
	if err != nil {
		return nil, errors.New("player not found")
	}

	// Try to get a match from the queue
	select {
	case match := <-s.queue:
		// A match was found, a player was waiting.
		// Add the new player and start the game.
		match.Players[1] = playerID
		err := s.gameRepo.Update(match.ID, match)
		if err != nil {
			// If we fail to update, something is wrong. We should probably try to requeue the original match.
			// For now, we'll just return the error.
			return nil, err
		}
		return match, nil
	default:
		// No match in queue, create a new one and queue it up.
		matchID := strconv.FormatUint(s.matchCounter, 10)
		s.matchCounter++

		match := &domain.Match{
			ID:      matchID,
			Players: [2]string{playerID, ""}, // Player 2 is empty for now
		}

		// Save the initial match to the repository
		err := s.gameRepo.Create(matchID, match)
		if err != nil {
			return nil, err
		}

		// Put the new match in the queue
		s.queue <- match
		// Return nil, as the match is not ready yet.
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
		// A user should always have one of each card type after buying a pack.
		// This error suggests they are trying to play a card type they don't have.
		return errors.New("player does not own a card of this type")
	}

	// Find player index and store move
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

	// Check if opponent has moved and determine winner
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

func (s *GameServiceImplementation) PLayerSurrender(gameID, player string) error {
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
