package data

import (
	"cod-server/internal/domain"
	"testing"
)

func TestInMemoryRepository(t *testing.T) {
	repo := NewInMemoryRepository[*domain.User]()

	t.Run("Create User", func(t *testing.T) {
		user := &domain.User{ID: "user1", Username: "testuser1"}
		err := repo.Create(user.ID, user)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		if _, exists := repo.entities.Get(user.ID); !exists {
			t.Errorf("User was not created in repository")
		}
	})

	t.Run("Read User", func(t *testing.T) {
		user, err := repo.Read("user1")
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
		if user.Username != "testuser1" {
			t.Errorf("Expected username 'testuser1', got '%s'", user.Username)
		}

		_, err = repo.Read("nonexistent")
		if err == nil {
			t.Errorf("Expected error for non-existent user, got none")
		}
	})

	t.Run("Update User", func(t *testing.T) {
		updatedUser := &domain.User{ID: "user1", Username: "updateduser"}
		err := repo.Update("user1", updatedUser)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}
		user, _ := repo.Read("user1")
		if user.Username != "updateduser" {
			t.Errorf("Expected username 'updateduser', got '%s'", user.Username)
		}

		err = repo.Update("nonexistent", updatedUser)
		if err == nil {
			t.Errorf("Expected error for non-existent user update, got none")
		}
	})

	t.Run("List Users", func(t *testing.T) {
		users, err := repo.List()
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(users) != 1 {
			t.Errorf("Expected 1 user in list, got %d", len(users))
		}
	})

	t.Run("Delete User", func(t *testing.T) {
		err := repo.Delete("user1")
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
		_, err = repo.Read("user1")
		if err == nil {
			t.Errorf("Expected error for deleted user, got none")
		}

		err = repo.Delete("nonexistent")
		if err == nil {
			t.Errorf("Expected error for deleting non-existent user, got none")
		}
	})
}
