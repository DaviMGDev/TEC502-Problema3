package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"cod-server/internal/utils"
	"log"
	"testing"
)

func TestGameService_StartGame(t *testing.T) {
	gameRepo := data.NewInMemoryRepository[*domain.Match]()
	userRepo := data.NewInMemoryRepository[*domain.User]()

	gameService := NewGameService(gameRepo, userRepo)

	userRepo.Create("p1", &domain.User{ID: "p1", Username: "player1", Cards: utils.NewSafeMap[string, *domain.Card]()})
	userRepo.Create("p2", &domain.User{ID: "p2", Username: "player2", Cards: utils.NewSafeMap[string, *domain.Card]()})

	t.Run("First player queues up", func(t *testing.T) {
		log.Printf("Running TestGameService_StartGame: First player queues up")

		match, err := gameService.StartGame("p1")
		if err != nil {
			t.Fatalf("Expected no error for first player, got %v", err)
		}
		if match != nil {
			t.Fatalf("Expected a nil match for the first player, got %+v", match)
		}

		if len(gameService.queue) != 1 {
			t.Errorf("Expected queue length to be 1, got %d", len(gameService.queue))
		}

		allMatches, _ := gameRepo.List()
		if len(allMatches) != 1 {
			t.Fatalf("Expected 1 match to be created in the repo, found %d", len(allMatches))
		}
		if allMatches[0].Players[0] != "p1" {
			t.Errorf("Expected player 1 of the queued match to be 'p1', got %s", allMatches[0].Players[0])
		}
		if allMatches[0].Players[1] != "" {
			t.Errorf("Expected player 2 of the queued match to be empty, but got %s", allMatches[0].Players[1])
		}
		log.Printf("TestGameService_StartGame 'First player queues up' passed.")
	})

	t.Run("Second player forms a match", func(t *testing.T) {
		log.Printf("Running TestGameService_StartGame: Second player forms a match")

		formedMatch, err := gameService.StartGame("p2")
		if err != nil {
			t.Fatalf("Expected no error for second player, got %v", err)
		}
		if formedMatch == nil {
			t.Fatal("Expected a non-nil match to be formed for the second player")
		}

		if len(gameService.queue) != 0 {
			t.Errorf("Expected queue to be empty after match formation, but length is %d", len(gameService.queue))
		}

		if formedMatch.Players[0] != "p1" || formedMatch.Players[1] != "p2" {
			t.Errorf("Expected formed match to have players p1 and p2, got %v", formedMatch.Players)
		}

		updatedMatch, err := gameRepo.Read(formedMatch.ID)
		if err != nil {
			t.Fatalf("Could not read updated match from repo: %v", err)
		}
		if updatedMatch.Players[1] != "p2" {
			t.Errorf("Expected player 2 to be updated in the repo, but got %s", updatedMatch.Players[1])
		}
		log.Printf("TestGameService_StartGame 'Second player forms a match' passed.")
	})
}

func TestGameService_MakeMove(t *testing.T) {

	setup := func() (*GameServiceImplementation, data.Repository[*domain.Match], data.Repository[*domain.User]) {
		gameRepo := data.NewInMemoryRepository[*domain.Match]()
		userRepo := data.NewInMemoryRepository[*domain.User]()
		gameService := NewGameService(gameRepo, userRepo)

		user1 := &domain.User{ID: "p1", Username: "player1", Cards: utils.NewSafeMap[string, *domain.Card]()}
		user1.Cards.Set("rock1", &domain.Card{ID: "rock1", Type: domain.Rock, Level: 5})
		user1.Cards.Set("paper1", &domain.Card{ID: "paper1", Type: domain.Paper, Level: 5})
		userRepo.Create("p1", user1)

		user2 := &domain.User{ID: "p2", Username: "player2", Cards: utils.NewSafeMap[string, *domain.Card]()}
		user2.Cards.Set("rock2", &domain.Card{ID: "rock2", Type: domain.Rock, Level: 5})
		user2.Cards.Set("scissors2", &domain.Card{ID: "scissors2", Type: domain.Scissors, Level: 5})
		userRepo.Create("p2", user2)

		gameRepo.Create("match1", &domain.Match{ID: "match1", Players: [2]string{"p1", "p2"}, Moves: [2]*domain.Card{}})

		return gameService, gameRepo, userRepo
	}

	t.Run("Successful First Move", func(t *testing.T) {
		gameService, gameRepo, _ := setup()
		log.Printf("Running TestGameService_MakeMove: Successful First Move")

		err := gameService.MakeMove("match1", "p1", "rock1")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		match, _ := gameRepo.Read("match1")
		if match.Moves[0] == nil {
			t.Fatal("Player 1's move was not recorded")
		}
		if match.Moves[0].ID != "rock1" {
			t.Errorf("Incorrect card recorded for player 1")
		}
		if match.Winner != "" {
			t.Errorf("Winner should not be decided yet, but got %s", match.Winner)
		}
		log.Printf("TestGameService_MakeMove 'Successful First Move' passed.")
	})

	t.Run("Second move determines winner", func(t *testing.T) {
		gameService, gameRepo, _ := setup()
		log.Printf("Running TestGameService_MakeMove: Second move determines winner")

		gameService.MakeMove("match1", "p1", "rock1")
		err := gameService.MakeMove("match1", "p2", "scissors2")
		if err != nil {
			t.Fatalf("Expected no error on second move, got %v", err)
		}

		match, _ := gameRepo.Read("match1")
		if match.Winner != "p1" {
			t.Errorf("Expected winner to be p1, got %s", match.Winner)
		}
		log.Printf("TestGameService_MakeMove 'Second move determines winner' passed.")
	})

	t.Run("Second move determines draw", func(t *testing.T) {
		gameService, gameRepo, _ := setup()
		log.Printf("Running TestGameService_MakeMove: Second move determines draw")

		gameService.MakeMove("match1", "p1", "rock1")
		err := gameService.MakeMove("match1", "p2", "rock2")
		if err != nil {
			t.Fatalf("Expected no error on second move, got %v", err)
		}

		match, _ := gameRepo.Read("match1")
		if match.Winner != "draw" {
			t.Errorf("Expected winner to be 'draw', got %s", match.Winner)
		}
		log.Printf("TestGameService_MakeMove 'Second move determines draw' passed.")
	})

	t.Run("Make Move on concluded match", func(t *testing.T) {
		gameService, gameRepo, _ := setup()
		log.Printf("Running TestGameService_MakeMove: Make Move on concluded match")

		match, _ := gameRepo.Read("match1")
		match.Winner = "p1"
		gameRepo.Update("match1", match)

		err := gameService.MakeMove("match1", "p2", "rock2")
		if err == nil {
			t.Errorf("Expected an error for making a move on a concluded match, got none")
		}
		log.Printf("TestGameService_MakeMove 'Make Move on concluded match' passed.")
	})

	t.Run("Make Move with invalid card", func(t *testing.T) {
		gameService, _, _ := setup()
		log.Printf("Running TestGameService_MakeMove: Make Move with invalid card")
		err := gameService.MakeMove("match1", "p1", "nonexistentcard")
		if err == nil {
			t.Errorf("Expected an error for using a card not owned by the player, got none")
		}
		log.Printf("TestGameService_MakeMove 'Make Move with invalid card' passed.")
	})
}

func TestGameService_GetGameState(t *testing.T) {
	gameRepo := data.NewInMemoryRepository[*domain.Match]()
	userRepo := data.NewInMemoryRepository[*domain.User]()
	gameService := NewGameService(gameRepo, userRepo)

	matchID := "match123"
	gameRepo.Create(matchID, &domain.Match{ID: matchID, Players: [2]string{"p1", "p2"}})

	t.Run("Successful Get Game State", func(t *testing.T) {
		log.Printf("Running TestGameService_GetGameState: Successful Get Game State")
		match, err := gameService.GetGameState(matchID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if match.ID != matchID {
			t.Errorf("Expected match ID %s, got %s", matchID, match.ID)
		}
		log.Printf("TestGameService_GetGameState 'Successful Get Game State' passed.")
	})

	t.Run("Get Game State Non-Existent Match", func(t *testing.T) {
		log.Printf("Running TestGameService_GetGameState: Get Game State Non-Existent Match")
		_, err := gameService.GetGameState("nonexistent")
		if err == nil {
			t.Errorf("Expected an error for non-existent match, got none")
		}
		log.Printf("TestGameService_GetGameState 'Get Game State Non-Existent Match' passed.")
	})
}

func TestGameService_PlayerSurrender(t *testing.T) {
	gameRepo := data.NewInMemoryRepository[*domain.Match]()
	userRepo := data.NewInMemoryRepository[*domain.User]()
	gameService := NewGameService(gameRepo, userRepo)

	matchID := "match123"
	player1ID := "player1"
	player2ID := "player2"
	gameRepo.Create(matchID, &domain.Match{ID: matchID, Players: [2]string{player1ID, player2ID}})

	t.Run("Successful Player Surrender", func(t *testing.T) {
		log.Printf("Running TestGameService_PlayerSurrender: Successful Player Surrender")
		err := gameService.PLayerSurrender(matchID, player1ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		match, _ := gameRepo.Read(matchID)
		if match.Winner != player2ID {
			t.Errorf("Expected winner to be %s, got %s", player2ID, match.Winner)
		}
		log.Printf("TestGameService_PlayerSurrender 'Successful Player Surrender' passed.")
	})

	t.Run("Player Surrender Non-Existent Match", func(t *testing.T) {
		log.Printf("Running TestGameService_PlayerSurrender: Player Surrender Non-Existent Match")
		err := gameService.PLayerSurrender("nonexistent", player1ID)
		if err == nil {
			t.Errorf("Expected error for non-existent match, got none")
		}
		log.Printf("TestGameService_PlayerSurrender 'Player Surrender Non-Existent Match' passed.")
	})

	t.Run("Player Surrender with player not in match", func(t *testing.T) {
		log.Printf("Running TestGameService_PlayerSurrender: Player Surrender with player not in match")
		err := gameService.PLayerSurrender(matchID, "player3")
		if err == nil {
			t.Errorf("Expected error for player not in match, got none")
		}
		log.Printf("TestGameService_PlayerSurrender 'Player Surrender with player not in match' passed.")
	})
}
