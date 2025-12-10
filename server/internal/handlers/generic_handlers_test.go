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
	// Spies
	RegisterCalled bool
	RegisterUsr    string
	RegisterPwd    string
	LoginCalled    bool
	LoginUsr       string
	LoginPwd       string

	// Stubs
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
	// Spies
	BuyCardPackCalled   bool
	BuyCardPackUserID   string
	TradeCalled         bool // Renamed from SwapCardCalled
	TradeUser1          string // Renamed from SwapCardUser1
	TradeUser2          string // Renamed from SwapCardUser2
	TradeCardType       string // Renamed from SwapCardCardType
	ListUserCardsCalled bool
	ListUserCardsUserID string

	// Stubs
	BuyCardPackReturnError    error
	TradeReturnError          error // Renamed from SwapCardReturnError
	ListUserCardsReturn       []*domain.Card
	ListUserCardsReturnError  error
}

func (m *MockCardService) BuyCardPack(userID string) error {
	m.BuyCardPackCalled = true
	m.BuyCardPackUserID = userID
	return m.BuyCardPackReturnError
}
func (m *MockCardService) Trade(user1, user2, cardType string) error { // Renamed from SwapCard
	m.TradeCalled = true
	m.TradeUser1 = user1
	m.TradeUser2 = user2
	m.TradeCardType = cardType
	return m.TradeReturnError
}
func (m *MockCardService) ListUserCards(userID string) ([]*domain.Card, error) {
	m.ListUserCardsCalled = true
	m.ListUserCardsUserID = userID
	return m.ListUserCardsReturn, m.ListUserCardsReturnError
}

type MockGameService struct {
	// Spies
	StartGameCalled       bool
	StartGamePlayerID     string
	MakeMoveCalled        bool
	MakeMoveGameID        string
	MakeMovePlayerID      string
	MakeMoveCardID        string
	GetGameStateCalled    bool
	GetGameStateGameID    string
	PlayerSurrenderCalled bool // Renamed from PLayerSurrenderCalled
	PlayerSurrenderGameID string // Renamed from PLayerSurrenderGameID
	PlayerSurrenderPlayer string // Renamed from PLayerSurrenderPlayer

	// Stubs
	StartGameReturnMatch        *domain.Match
	StartGameReturnError error
	MakeMoveReturnError  error
	GetGameStateReturnMatch *domain.Match
	GetGameStateReturnError error
	PlayerSurrenderReturnError error // Renamed from PLayerSurrenderReturnError
}

func (m *MockGameService) StartGame(playerID string) (*domain.Match, error) {
	m.StartGameCalled = true
	m.StartGamePlayerID = playerID
	return m.StartGameReturnMatch, m.StartGameReturnError
}
func (m *MockGameService) MakeMove(gameID, player, cardID string) error {
	m.MakeMoveCalled = true
	m.MakeMoveGameID = gameID
	m.MakeMovePlayerID = player
	m.MakeMoveCardID = cardID
	return m.MakeMoveReturnError
}
func (m *MockGameService) GetGameState(gameID string) (*domain.Match, error) {
	m.GetGameStateCalled = true
	m.GetGameStateGameID = gameID
	return m.GetGameStateReturnMatch, m.GetGameStateReturnError
}
func (m *MockGameService) PlayerSurrender(gameID, player string) error { // Renamed from PLayerSurrender
	m.PlayerSurrenderCalled = true
	m.PlayerSurrenderGameID = gameID
	m.PlayerSurrenderPlayer = player
	return m.PlayerSurrenderReturnError
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
			Payload: utils.Dict{"username": "testuser", "password": "testpassword", "client_id": "client123"},
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
		if userID, ok := resp.Payload["user_id"].(string); !ok || userID == "" { // Check if user_id is populated
			t.Errorf("Expected non-empty user_id in payload, got '%v'", resp.Payload["user_id"])
		}
		log.Println("TestOnRegisterEvent 'calls user service with correct params and returns success' completed (expected to fail).")
	})

	t.Run("returns error on missing payload fields", func(t *testing.T) {
		clients := setupTestHandlers()

		resp := clients.Handlers.OnRegisterEvent(protocol.Event{
			Method:  "register",
			Payload: utils.Dict{"username": "testuser", "client_id": "client123"}, // Missing password
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
			Payload: utils.Dict{"username": "erroruser", "password": "password", "client_id": "client123"},
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

func TestOnLoginEvent(t *testing.T) {
	t.Run("Successful Login calls service with correct params and returns success", func(t *testing.T) {
		clients := setupTestHandlers()
		clients.UserService.LoginReturnID = "existingUserID"

		inputEvent := protocol.Event{
			Method:  "login",
			Payload: utils.Dict{"username": "testuser", "password": "testpassword", "client_id": "client123"},
		}

		resp := clients.Handlers.OnLoginEvent(inputEvent)

		if !clients.UserService.LoginCalled {
			t.Error("Expected Login to be called, but it was not.")
		}
		if clients.UserService.LoginUsr != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", clients.UserService.LoginUsr)
		}
		if clients.UserService.LoginPwd != "testpassword" {
			t.Errorf("Expected password 'testpassword', got '%s'", clients.UserService.LoginPwd)
		}

		if resp.Method != "login" {
			t.Errorf("Expected response method 'login', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "success" {
			t.Errorf("Expected status 'success', got '%v'", resp.Payload["status"])
		}
		if userID, ok := resp.Payload["user_id"].(string); !ok || userID != "existingUserID" {
			t.Errorf("Expected user_id 'existingUserID' in payload, got '%v'", resp.Payload["user_id"])
		}
		log.Println("TestOnLoginEvent 'Successful Login' completed (expected to fail).")
	})

	t.Run("returns error on missing payload fields", func(t *testing.T) {
		clients := setupTestHandlers()
		
		resp := clients.Handlers.OnLoginEvent(protocol.Event{
			Method:  "login",
			Payload: utils.Dict{"username": "testuser", "client_id": "client123"}, // Missing password
		})

		if resp.Method != "login" {
			t.Errorf("Expected response method 'login', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg == "" {
			t.Error("Expected an error message in payload, but got nil or empty")
		}
		log.Println("TestOnLoginEvent 'returns error on missing payload fields' completed (expected to fail).")
	})

	t.Run("returns error on service failure", func(t *testing.T) {
		clients := setupTestHandlers()
		clients.UserService.LoginReturnError = errors.New("invalid credentials")
		
		resp := clients.Handlers.OnLoginEvent(protocol.Event{
			Method:  "login",
			Payload: utils.Dict{"username": "erroruser", "password": "password", "client_id": "client123"},
		})

		if resp.Method != "login" {
			t.Errorf("Expected response method 'login', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg == "" {
			t.Errorf("Expected error message 'invalid credentials', got '%v'", resp.Payload["status_message"])
		}
		log.Println("TestOnLoginEvent 'returns error on service failure' completed (expected to fail).")
	})
}

func TestOnBuyCardPackEvent(t *testing.T) {
	t.Run("calls card service with correct params and returns success", func(t *testing.T) {
		clients := setupTestHandlers()
		
		inputEvent := protocol.Event{
			Method:  "buy_pack",
			Payload: utils.Dict{"user_id": "user123", "client_id": "clientABC"},
		}
		
		resp := clients.Handlers.OnBuyCardPackEvent(inputEvent)

		if !clients.CardService.BuyCardPackCalled {
			t.Error("Expected BuyCardPack to be called, but it was not.")
		}
		if clients.CardService.BuyCardPackUserID != "user123" {
			t.Errorf("Expected user_id 'user123', got '%s'", clients.CardService.BuyCardPackUserID)
		}
		
		if resp.Method != "buy_pack" {
			t.Errorf("Expected response method 'buy_pack', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "success" {
			t.Errorf("Expected status 'success', got '%v'", resp.Payload["status"])
		}
		log.Println("TestOnBuyCardPackEvent 'calls card service with correct params and returns success' completed (expected to fail).")
	})

	t.Run("returns error on missing payload fields", func(t *testing.T) {
		clients := setupTestHandlers()
		
		resp := clients.Handlers.OnBuyCardPackEvent(protocol.Event{
			Method:  "buy_pack",
			Payload: utils.Dict{"client_id": "clientABC"}, // Missing user_id
		})

		if resp.Method != "buy_pack" {
			t.Errorf("Expected response method 'buy_pack', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg == "" {
			t.Error("Expected an error message in payload, but got nil or empty")
		}
		log.Println("TestOnBuyCardPackEvent 'returns error on missing payload fields' completed (expected to fail).")
	})

	t.Run("returns error on service failure", func(t *testing.T) {
		clients := setupTestHandlers()
		clients.CardService.BuyCardPackReturnError = errors.New("buy pack failed")
		
		resp := clients.Handlers.OnBuyCardPackEvent(protocol.Event{
			Method:  "buy_pack",
			Payload: utils.Dict{"user_id": "erroruser", "client_id": "clientABC"},
		})

		if resp.Method != "buy_pack" {
			t.Errorf("Expected response method 'buy_pack', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg != "buy pack failed" {
			t.Errorf("Expected error message 'buy pack failed', got '%v'", resp.Payload["status_message"])
		}
		log.Println("TestOnBuyCardPackEvent 'returns error on service failure' completed (expected to fail).")
	})
}

func TestOnTradeEvent(t *testing.T) { // Renamed from TestOnSwapCardEvent
	t.Run("calls card service with correct params and returns success", func(t *testing.T) {
		clients := setupTestHandlers()
		
		inputEvent := protocol.Event{
			Method:  "trade", 
			Payload: utils.Dict{"user_id": "user1", "card_type": "rock", "target_user_id": "user2", "client_id": "clientABC"}, // Updated payload fields
		}
		
		resp := clients.Handlers.OnTradeEvent(inputEvent) // Renamed from OnSwapCardEvent

		if !clients.CardService.TradeCalled { // Renamed from SwapCardCalled
			t.Error("Expected Trade to be called, but it was not.")
		}
		if clients.CardService.TradeUser1 != "user1" {
			t.Errorf("Expected user1 'user1', got '%s'", clients.CardService.TradeUser1)
		}
		if clients.CardService.TradeUser2 != "user2" {
			t.Errorf("Expected user2 'user2', got '%s'", clients.CardService.TradeUser2)
		}
		if clients.CardService.TradeCardType != "rock" {
			t.Errorf("Expected card_type 'rock', got '%s'", clients.CardService.TradeCardType)
		}
		
		if resp.Method != "trade" { 
			t.Errorf("Expected response method 'trade', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "success" {
			t.Errorf("Expected status 'success', got '%v'", resp.Payload["status"])
		}
		log.Println("TestOnTradeEvent 'calls card service with correct params and returns success' completed (expected to fail).")
	})

	t.Run("returns error on missing payload fields", func(t *testing.T) {
		clients := setupTestHandlers()
		
		resp := clients.Handlers.OnTradeEvent(protocol.Event{ // Renamed from OnSwapCardEvent
			Method:  "trade",
			Payload: utils.Dict{"user_id": "user1", "client_id": "clientABC"}, // Missing card_type and target_user_id
		})

		if resp.Method != "trade" {
			t.Errorf("Expected response method 'trade', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg == "" {
			t.Error("Expected an error message in payload, but got nil or empty")
		}
		log.Println("TestOnTradeEvent 'returns error on missing payload fields' completed (expected to fail).")
	})

	t.Run("returns error on service failure", func(t *testing.T) {
		clients := setupTestHandlers()
		clients.CardService.TradeReturnError = errors.New("trade failed") // Renamed from SwapCardReturnError
		
		resp := clients.Handlers.OnTradeEvent(protocol.Event{ // Renamed from OnSwapCardEvent
			Method:  "trade",
			Payload: utils.Dict{"user_id": "user1", "card_type": "rock", "target_user_id": "user2", "client_id": "clientABC"},
		})

		if resp.Method != "trade" {
			t.Errorf("Expected response method 'trade', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg != "trade failed" {
			t.Errorf("Expected error message 'trade failed', got '%v'", resp.Payload["status_message"])
		}
		log.Println("TestOnTradeEvent 'returns error on service failure' completed (expected to fail).")
	})
}

func TestOnListUserCardsEvent(t *testing.T) {
	t.Run("calls card service with correct params and returns success", func(t *testing.T) {
		clients := setupTestHandlers()
		clients.CardService.ListUserCardsReturn = []*domain.Card{
			{ID: "c1", Type: domain.Rock, Level: 1},
			{ID: "c2", Type: domain.Paper, Level: 5},
		}

		inputEvent := protocol.Event{
			Method:  "list_cards",
			Payload: utils.Dict{"user_id": "user123", "client_id": "clientABC"},
		}
		
		resp := clients.Handlers.OnListUserCardsEvent(inputEvent)

		if !clients.CardService.ListUserCardsCalled {
			t.Error("Expected ListUserCards to be called, but it was not.")
		}
		
		if resp.Method != "list_cards" {
			t.Errorf("Expected response method 'list_cards', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "success" {
			t.Errorf("Expected status 'success', got '%v'", resp.Payload["status"])
		}
		if cards, ok := resp.Payload["cards"].([]map[string]interface{}); !ok || len(cards) != 2 {
			t.Errorf("Expected 2 cards in payload, got %v", resp.Payload["cards"])
		}
		log.Println("TestOnListUserCardsEvent 'calls card service with correct params and returns success' completed (expected to fail).")
	})

	t.Run("returns error on missing payload fields", func(t *testing.T) {
		clients := setupTestHandlers()
		
		resp := clients.Handlers.OnListUserCardsEvent(protocol.Event{
			Method:  "list_cards",
			Payload: utils.Dict{"client_id": "clientABC"}, // Missing user_id
		})

		if resp.Method != "list_cards" {
			t.Errorf("Expected response method 'list_cards', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg == "" {
			t.Error("Expected an error message in payload, but got nil or empty")
		}
		log.Println("TestOnListUserCardsEvent 'returns error on missing payload fields' completed (expected to fail).")
	})

	t.Run("returns error on service failure", func(t *testing.T) {
		clients := setupTestHandlers()
		clients.CardService.ListUserCardsReturnError = errors.New("list cards failed")
		
		resp := clients.Handlers.OnListUserCardsEvent(protocol.Event{
			Method:  "list_cards",
			Payload: utils.Dict{"user_id": "erroruser", "client_id": "clientABC"},
		})

		if resp.Method != "list_cards" {
			t.Errorf("Expected response method 'list_cards', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg == "" {
			t.Error("Expected an error message in payload, but got nil or empty")
		}
		log.Println("TestOnListUserCardsEvent 'returns error on service failure' completed (expected to fail).")
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
		log.Println("TestOnStartMatchEvent 'calls game service with correct params and returns success' completed (expected to fail).")
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
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg == "" {
			t.Errorf("Expected error message 'game service error', got '%v'", resp.Payload["status_message"])
		}
		log.Println("TestOnStartMatchEvent 'returns error on service failure' completed (expected to fail).")
	})
}

func TestOnPlayerSurrenderEvent(t *testing.T) { // Renamed from TestOnPLayerSurrenderEvent
	t.Run("calls game service with correct params and returns success", func(t *testing.T) {
		clients := setupTestHandlers()
		clients.GameService.PlayerSurrenderReturnError = nil // Renamed from PLayerSurrenderReturnError

		inputEvent := protocol.Event{
			Method:  "surrender",
			Payload: utils.Dict{"user_id": "player1", "match_id": "match123", "client_id": "clientABC"},
		}

		resp := clients.Handlers.OnPlayerSurrenderEvent(inputEvent) // Renamed from OnPLayerSurrenderEvent

		if !clients.GameService.PlayerSurrenderCalled { // Renamed from PLayerSurrenderCalled
			t.Error("Expected PlayerSurrender to be called, but it was not.")
		}
		if clients.GameService.PlayerSurrenderGameID != "match123" { // Renamed from PLayerSurrenderGameID
			t.Errorf("Expected gameID 'match123', got '%s'", clients.GameService.PlayerSurrenderGameID)
		}
		if clients.GameService.PlayerSurrenderPlayer != "player1" { // Renamed from PLayerSurrenderPlayer
			t.Errorf("Expected player 'player1', got '%s'", clients.GameService.PlayerSurrenderPlayer)
		}
		
		if resp.Method != "surrender" {
			t.Errorf("Expected response method 'surrender', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "success" {
			t.Errorf("Expected status 'success', got '%v'", resp.Payload["status"])
		}
		log.Println("TestOnPlayerSurrenderEvent 'calls game service with correct params and returns success' completed (expected to fail).")
	})

	t.Run("returns error on missing payload fields", func(t *testing.T) {
		clients := setupTestHandlers()
		
		resp := clients.Handlers.OnPlayerSurrenderEvent(protocol.Event{ // Renamed from OnPLayerSurrenderEvent
			Method:  "surrender",
			Payload: utils.Dict{"user_id": "player1", "client_id": "clientABC"}, // Missing match_id
		})

		if resp.Method != "surrender" {
			t.Errorf("Expected response method 'surrender', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg == "" {
			t.Error("Expected an error message in payload, but got nil or empty")
		}
		log.Println("TestOnPlayerSurrenderEvent 'returns error on missing payload fields' completed (expected to fail).")
	})

	t.Run("returns error on service failure", func(t *testing.T) {
		clients := setupTestHandlers()
		clients.GameService.PlayerSurrenderReturnError = errors.New("surrender failed") // Renamed from PLayerSurrenderReturnError
		
		resp := clients.Handlers.OnPlayerSurrenderEvent(protocol.Event{ // Renamed from OnPLayerSurrenderEvent
			Method:  "surrender",
			Payload: utils.Dict{"user_id": "player1", "match_id": "match123", "client_id": "clientABC"},
		})

		if resp.Method != "surrender" {
			t.Errorf("Expected response method 'surrender', got '%s'", resp.Method)
		}
		if status, ok := resp.Payload["status"].(string); !ok || status != "error" {
			t.Errorf("Expected status 'error', got '%v'", resp.Payload["status"])
		}
		if msg, ok := resp.Payload["status_message"].(string); !ok || msg != "surrender failed" {
			t.Errorf("Expected error message 'surrender failed', got '%v'", resp.Payload["status_message"])
		}
		log.Println("TestOnPlayerSurrenderEvent 'returns error on service failure' completed (expected to fail).")
	})
}