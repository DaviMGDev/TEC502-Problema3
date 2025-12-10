package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"cod-server/internal/utils"
	"log"
	"testing"
)

// --- UserService Tests ---

func TestUserService_Register(t *testing.T) {
	t.Run("Successful Registration", func(t *testing.T) {
		log.Printf("Running TestUserService_Register: Successful Registration")
		userRepo := data.NewInMemoryRepository[*domain.User]()
		userService := NewUserService(userRepo)
		
		createdUser, err := userService.Register("testuser", "testpass")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if createdUser == nil {
			t.Fatal("Expected a non-nil user to be returned")
		}
		if createdUser.Username != "testuser" {
			t.Errorf("Expected created user to have username 'testuser', got '%s'", createdUser.Username)
		}
		
		user, err := userRepo.Read("testuser")
		if err != nil {
			t.Fatalf("Expected user to be created in repo, got error %v", err)
		}
		if user.Username != "testuser" {
			t.Errorf("Expected username 'testuser' in repo, got '%s'", user.Username)
		}
		log.Printf("TestUserService_Register 'Successful Registration' passed.")
	})

	t.Run("Register existing user", func(t *testing.T) {
		log.Printf("Running TestUserService_Register: Register existing user")
		userRepo := data.NewInMemoryRepository[*domain.User]()
		userService := NewUserService(userRepo)
		
		// First registration should succeed
		_, err1 := userService.Register("testuser", "pass1")
		if err1 != nil {
			t.Fatalf("First registration failed unexpectedly: %v", err1)
		}
		
		// Second registration should fail
		_, err2 := userService.Register("testuser", "pass2")
		if err2 == nil {
			t.Fatalf("Expected error for existing user, got none")
		}
		log.Printf("TestUserService_Register 'Register existing user' passed.")
	})
}

func TestUserService_Login(t *testing.T) {
	mockRepo := data.NewInMemoryRepository[*domain.User]()
	userService := NewUserService(mockRepo)

	// Pre-register a user
	_, _ = userService.Register("loginuser", "loginpass")

	t.Run("Successful Login", func(t *testing.T) {
		log.Printf("Running TestUserService_Login: Successful Login")
		userID, err := userService.Login("loginuser", "loginpass")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if userID == "" {
			t.Errorf("Expected a non-empty user ID")
		}
		log.Printf("TestUserService_Login 'Successful Login' passed.")
	})

	t.Run("Login with wrong password", func(t *testing.T) {
		log.Printf("Running TestUserService_Login: Login with wrong password")
		_, err := userService.Login("loginuser", "wrongpass")
		if err == nil {
			t.Fatalf("Expected error for wrong password, got none")
		}
		if err.Error() != "invalid credentials" {
			t.Errorf("Expected 'invalid credentials' error, got %v", err)
		}
		log.Printf("TestUserService_Login 'Login with wrong password' passed.")
	})

	t.Run("Login non-existent user", func(t *testing.T) {
		log.Printf("Running TestUserService_Login: Login non-existent user")
		_, err := userService.Login("nonexistent", "pass")
		if err == nil {
			t.Fatalf("Expected error for non-existent user, got none")
		}
		log.Printf("TestUserService_Login 'Login non-existent user' passed.")
	})
}

// --- CardService Tests ---

func TestCardService_BuyCardPack(t *testing.T) {
	mockRepo := data.NewInMemoryRepository[*domain.User]()
	userService := NewUserService(mockRepo)
	_, _ = userService.Register("buyeruser", "pass")

	t.Run("Successful Buy Card Pack", func(t *testing.T) {
		// NewCardService now correctly initializes and stocks the packages list.
		cardService := NewCardService(mockRepo)
		log.Printf("Running TestCardService_BuyCardPack: Successful Buy Card Pack")
		
		initialPackCount := cardService.packages.Size()
		if initialPackCount == 0 {
			t.Fatalf("Expected initial packages to be stocked, but count is 0")
		}

		err := cardService.BuyCardPack("buyeruser")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cardService.packages.Size() != initialPackCount-1 {
			t.Errorf("Expected package count to decrease by 1, got %d", cardService.packages.Size())
		}
		user, _ := mockRepo.Read("buyeruser")
		if user.Cards.Size() != 3 { // Rock, Paper, Scissors
			t.Errorf("Expected user to have 3 cards, got %d", user.Cards.Size())
		}
		log.Printf("TestCardService_BuyCardPack 'Successful Buy Card Pack' passed.")
	})

	t.Run("Buy Card Pack non-existent user", func(t *testing.T) {
		cardService := NewCardService(mockRepo)
		log.Printf("Running TestCardService_BuyCardPack: Buy Card Pack non-existent user")
		err := cardService.BuyCardPack("nonexistent")
		if err == nil {
			t.Fatalf("Expected error for non-existent user, got none")
		}
		log.Printf("TestCardService_BuyCardPack 'Buy Card Pack non-existent user' passed.")
	})

	t.Run("checkStock replenishes packages", func(t *testing.T) {
		log.Printf("Running TestCardService_BuyCardPack: checkStock replenishes packages")
		cardService := NewCardService(mockRepo)
		
		// Manually drain the list to test replenishment logic in isolation
		for cardService.packages.Size() > 0 {
			cardService.packages.Pop()
		}
		// Add one package back to be able to Pop it
		cardService.packages.Append(&domain.Package{
			ID:    "last_pack",
			Cards: utils.NewSafeMap[string, *domain.Card](),
		})

		if cardService.packages.Size() != 1 {
			t.Fatalf("Precondition failed: expected 1 package, got %d", cardService.packages.Size())
		}

		err := cardService.BuyCardPack("buyeruser")
		if err != nil {
			t.Fatalf("Expected no error on buy, got %v", err)
		}

		expectedStock := 64
		if cardService.packages.Size() != expectedStock {
			t.Errorf("Expected packages to be replenished to %d, but size is %d", expectedStock, cardService.packages.Size())
		}
		log.Printf("TestCardService_BuyCardPack 'checkStock replenishes packages' passed.")
	})
}

func TestCardService_Trade(t *testing.T) {
	mockRepo := data.NewInMemoryRepository[*domain.User]()
	cardService := NewCardService(mockRepo)
	userService := NewUserService(mockRepo)

	_, _ = userService.Register("userA", "passA")
	_, _ = userService.Register("userB", "passB")

	userA, _ := mockRepo.Read("userA")
	userB, _ := mockRepo.Read("userB")

	userA.Cards.Set(domain.Rock, &domain.Card{ID: "rockA", Type: domain.Rock, Level: 5})
	userB.Cards.Set(domain.Rock, &domain.Card{ID: "rockB", Type: domain.Rock, Level: 10})
	mockRepo.Update("userA", userA)
	mockRepo.Update("userB", userB)

	t.Run("Successful Trade", func(t *testing.T) {
		log.Printf("Running TestCardService_Trade: Successful Trade")
		err := cardService.Trade("userA", "userB", domain.Rock)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		updatedUserA, _ := mockRepo.Read("userA")
		updatedUserB, _ := mockRepo.Read("userB")

		cardA, _ := updatedUserA.Cards.Get(domain.Rock)
		cardB, _ := updatedUserB.Cards.Get(domain.Rock)

		if cardA.ID != "rockB" {
			t.Errorf("UserA should have rockB, got %s", cardA.ID)
		}
		if cardB.ID != "rockA" {
			t.Errorf("UserB should have rockA, got %s", cardB.ID)
		}
		log.Printf("TestCardService_Trade 'Successful Trade' passed.")
	})

	t.Run("Trade non-existent user", func(t *testing.T) {
		log.Printf("Running TestCardService_Trade: Trade non-existent user")
		err := cardService.Trade("nonexistent", "userA", domain.Rock)
		if err == nil {
			t.Fatalf("Expected error for non-existent user, got none")
		}
		log.Printf("TestCardService_Trade 'Trade non-existent user' passed.")
	})

	t.Run("Trade Card type not owned", func(t *testing.T) {
		log.Printf("Running TestCardService_Trade: Trade Card type not owned")
		err := cardService.Trade("userA", "userB", domain.Scissors)
		if err == nil { // Changed from err != nil to err == nil because the service now returns an error
			t.Fatalf("Expected error when card type not owned, got none")
		}
		log.Printf("TestCardService_Trade 'Trade Card type not owned' passed.")
	})
}

func TestCardService_ListUserCards(t *testing.T) {
	mockRepo := data.NewInMemoryRepository[*domain.User]()
	cardService := NewCardService(mockRepo)
	userService := NewUserService(mockRepo)

	_, _ = userService.Register("listeruser", "pass")
	user, _ := mockRepo.Read("listeruser")
	user.Cards.Set(domain.Rock, &domain.Card{ID: "rockL", Type: domain.Rock, Level: 1})
	mockRepo.Update("listeruser", user)

	t.Run("Successful List User Cards", func(t *testing.T) {
		log.Printf("Running TestCardService_ListUserCards: Successful List User Cards")
		cards, err := cardService.ListUserCards("listeruser")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(cards) != 1 {
			t.Errorf("Expected 1 card, got %d", len(cards))
		}
		log.Printf("TestCardService_ListUserCards 'Successful List User Cards' passed.")
	})

	t.Run("List Cards non-existent user", func(t *testing.T) {
		log.Printf("Running TestCardService_ListUserCards: List Cards non-existent user")
		_, err := cardService.ListUserCards("nonexistent")
		if err == nil {
			t.Fatalf("Expected error for non-existent user, got none")
		}
		log.Printf("TestCardService_ListUserCards 'List Cards non-existent user' passed.")
	})
}