package handlers

import (
	"cod-server/internal/api/protocol"
	"cod-server/internal/domain"
	"cod-server/internal/utils"
	"errors"
	"log"
	"testing"
)

// --- Mocks (Spies) ---

type MockUserService struct {
	RegisterCalled bool
	RegisterUsr    string
	RegisterPwd    string
	LoginCalled    bool
	LoginUsr       string
	LoginPwd       string

	RegisterReturnUser  *domain.User
	RegisterReturnError error
	LoginReturnID       string
	LoginReturnError    error
}

func (m *MockUserService) Register(username, password string) (*domain.User, error) {
	m.RegisterCalled = true
	m.RegisterUsr = username
	m.RegisterPwd = password
	return m.RegisterReturnUser, m.RegisterReturnError
}
func (m *MockUserService) Login(username, password string) (string, error) {
	m.LoginCalled = true
	m.LoginUsr = username
	m.LoginPwd = password
	return m.LoginReturnID, m.LoginReturnError
}

type MockCardService struct {
	BuyCardPackCalled bool
	BuyCardPackUserID string
	SwapCardCalled    bool
	ListUserCardsCalled bool

	BuyCardPackReturnError error
	SwapCardReturnError    error
	ListUserCardsReturn    []*domain.Card
	ListUserCardsReturnError error
}

func (m *MockCardService) BuyCardPack(userID string) error {
	m.BuyCardPackCalled = true
	m.BuyCardPackUserID = userID
	return m.BuyCardPackReturnError
}
func (m *MockCardService) SwapCard(user1, user2, cardType string) error {
	m.SwapCardCalled = true
	return m.SwapCardReturnError
}
func (m *MockCardService) ListUserCards(userID string) ([]*domain.Card, error) {
	m.ListUserCardsCalled = true
	return m.ListUserCardsReturn, m.ListUserCardsReturnError
}

type MockGameService struct {
	StartGameCalled       bool
	StartGamePlayerID     string
	MakeMoveCalled        bool
	GetGameStateCalled    bool
	PLayerSurrenderCalled bool

	StartGameReturnMatch *domain.Match
	StartGameReturnError error
	MakeMoveReturnError  error
	GetGameStateReturnMatch *domain.Match
	GetGameStateReturnError error
	PLayerSurrenderReturnError error
}

func (m *MockGameService) StartGame(playerID string) (*domain.Match, error) {
	m.StartGameCalled = true
	m.StartGamePlayerID = playerID
	return m.StartGameReturnMatch, m.StartGameReturnError
}
func (m *MockGameService) MakeMove(gameID, player, cardID string) error {
	m.MakeMoveCalled = true
	return m.MakeMoveReturnError
}
func (m *MockGameService) GetGameState(gameID string) (*domain.Match, error) {
	m.GetGameStateCalled = true
	return m.GetGameStateReturnMatch, m.GetGameStateReturnError
}
func (m *MockGameService) PLayerSurrender(gameID, player string) error {
	m.PLayerSurrenderCalled = true
	return m.PLayerSurrenderReturnError
}

// --- Setup Helper ---

type TestHandlerClients struct {
	Handlers    *HandlersImplementation
	UserService *MockUserService
	CardService *MockCardService
	GameService *MockGameService
}

func setupTestHandlers() TestHandlerClients {
	mockUser := &MockUserService{}
	mockCard := &MockCardService{}
	mockGame := &MockGameService{}
	handlers := NewHandlers(mockUser, mockCard, mockGame)
	return TestHandlerClients{
		Handlers:    handlers,
		UserService: mockUser,
		CardService: mockCard,
		GameService: mockGame,
	}
}

// --- Handler Tests ---

func TestOnRegisterEvent(t *testing.T) {
	t.Run("calls user service with correct params and returns success", func(t *testing.T) {
		clients := setupTestHandlers()
		// Mock service to return a user
		clients.UserService.RegisterReturnUser = &domain.User{ID: "newUserID"}
		
		inputEvent := protocol.Event{
			Method:  "register",
			Payload: utils.Dict{"username": "testuser", "password": "testpassword"},
		}
		
		resp := clients.Handlers.OnRegisterEvent(inputEvent)

		if !clients.UserService.RegisterCalled {
			t.Error("Expected Register to be called, but it was not.")
		}
		if clients.UserService.RegisterUsr != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", clients.UserService.RegisterUsr)
		}
		if clients.UserService.RegisterPwd != "testpassword" {
			t.Errorf("Expected password 'testpassword', got '%s'", clients.UserService.RegisterPwd)
		}
		
		if resp.Method != "register" {
			t.Errorf("Expected response method 'register', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "success" {
			t.Errorf("Expected status 'success', got '%v'", resp.Payload["status"])
		}
		if userID, ok := resp.Payload["user_id"].(string); !ok || userID != "newUserID" {
			t.Errorf("Expected user_id 'newUserID' in payload, got '%v'", resp.Payload["user_id"])
		}
		log.Println("TestOnRegisterEvent 'calls user service with correct params and returns success' completed (expected to fail).")
	})

	t.Run("returns error on missing payload fields", func(t *testing.T) {
		clients := setupTestHandlers()
		
		resp := clients.Handlers.OnRegisterEvent(protocol.Event{
			Method:  "register",
			Payload: utils.Dict{"username": "testuser"}, // Missing password
		})

		if resp.Method != "register" {
			t.Errorf("Expected response method 'register', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg == "" {
			t.Error("Expected an error message in payload, but got nil or empty")
		}
		log.Println("TestOnRegisterEvent 'returns error on missing payload fields' completed (expected to fail).")
	})

	t.Run("returns error on service failure", func(t *testing.T) {
		clients := setupTestHandlers()
		clients.UserService.RegisterReturnError = errors.New("user service error")
		
		resp := clients.Handlers.OnRegisterEvent(protocol.Event{
			Method:  "register",
			Payload: utils.Dict{"username": "erroruser", "password": "password"},
		})

		if resp.Method != "register" {
			t.Errorf("Expected response method 'register', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg != "user service error" {
			t.Errorf("Expected error message 'user service error', got '%v'", resp.Payload["status_message"])
		}
		log.Println("TestOnRegisterEvent 'returns error on service failure' completed (expected to fail).")
	})
}

func TestOnStartMatchEvent(t *testing.T) {
	t.Run("calls game service with correct playerID for queuing", func(t *testing.T) {
		clients := setupTestHandlers()
		
		// Simulate StartGame returning nil match (player queued)
		clients.GameService.StartGameReturnMatch = nil 
		clients.GameService.StartGameReturnError = nil

		inputEvent := protocol.Event{
			Method:  "start_game",
			Payload: utils.Dict{"user_id": "p1"},
		}

		resp := clients.Handlers.OnStartMatchEvent(inputEvent)

		if !clients.GameService.StartGameCalled {
			t.Error("Expected StartGame to be called, but it was not.")
		}
		if clients.GameService.StartGamePlayerID != "p1" {
			t.Errorf("Expected StartGame to be called with 'p1', but got '%s'", clients.GameService.StartGamePlayerID)
		}
		
		if resp.Method != "start_game" {
			t.Errorf("Expected response method 'start_game', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "success" {
			t.Errorf("Expected status 'success', got '%s'", status)
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg != "player queued" {
			t.Errorf("Expected status_message 'player queued', got '%s'", msg)
		}
		if resp.Payload["match_id"] != nil { // Should not have a match ID if queued
			t.Errorf("Expected match_id to be nil, got '%v'", resp.Payload["match_id"])
		}
		log.Println("TestOnStartMatchEvent 'calls game service with correct playerID for queuing' completed (expected to fail).")
	})

	t.Run("forms match and returns match_id", func(t *testing.T) {
		clients := setupTestHandlers()
		matchResult := &domain.Match{ID: "match123", Players: [2]string{"p_waiting", "p2"}}
		clients.GameService.StartGameReturnMatch = matchResult
		clients.GameService.StartGameReturnError = nil

		resp := clients.Handlers.OnStartMatchEvent(protocol.Event{
			Method:  "start_game",
			Payload: utils.Dict{"user_id": "p2"},
		})

		if !clients.GameService.StartGameCalled {
			t.Error("Expected StartGame to be called, but it was not.")
		}
		if clients.GameService.StartGamePlayerID != "p2" {
			t.Errorf("Expected StartGame to be called with 'p2', but got '%s'", clients.GameService.StartGamePlayerID)
		}

		if resp.Method != "start_game" {
			t.Errorf("Expected response method 'start_game', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "success" {
			t.Errorf("Expected status 'success', got '%s'", status)
		}
		if matchID, ok := resp.Payload["match_id"].(string); !ok || matchID != "match123" {
			t.Errorf("Expected match_id 'match123', got '%v'", resp.Payload["match_id"])
		}
		log.Println("TestOnStartMatchEvent 'forms match and returns match_id' completed (expected to fail).")
	})

	t.Run("returns error on missing payload fields", func(t *testing.T) {
		clients := setupTestHandlers()
		
		resp := clients.Handlers.OnStartMatchEvent(protocol.Event{
			Method:  "start_game",
			Payload: utils.Dict{}, // Missing user_id
		})

		if resp.Method != "start_game" {
			t.Errorf("Expected response method 'start_game', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg == "" {
			t.Error("Expected an error message in payload, but got nil or empty")
		}
		log.Println("TestOnStartMatchEvent 'returns error on missing payload fields' completed (expected to fail).")
	})

	t.Run("returns error on service failure", func(t *testing.T) {
		clients := setupTestHandlers()
		clients.GameService.StartGameReturnError = errors.New("game service error")
		
		resp := clients.Handlers.OnStartMatchEvent(protocol.Event{
			Method:  "start_game",
			Payload: utils.Dict{"user_id": "p1"},
		})

		if resp.Method != "start_game" {
			t.Errorf("Expected response method 'start_game', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg != "game service error" {
			t.Errorf("Expected error message 'game service error', got '%v'", resp.Payload["status_message"])
		}
		log.Println("TestOnStartMatchEvent 'returns error on service failure' completed (expected to fail).")
	})
}