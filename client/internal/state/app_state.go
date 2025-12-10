package state

import (
	"sync"
)

// AppState holds the application state with thread-safe access
type AppState struct {
	mutex    sync.RWMutex
	userID   string
	username string
	matchID  string
}

// NewAppState creates a new application state instance
func NewAppState() *AppState {
	return &AppState{}
}

// GetUserID returns the current user ID
func (s *AppState) GetUserID() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.userID
}

// SetUserID sets the current user ID
func (s *AppState) SetUserID(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.userID = userID
}

// GetUsername returns the current username
func (s *AppState) GetUsername() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.username
}

// SetUsername sets the current username
func (s *AppState) SetUsername(username string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.username = username
}

// GetMatchID returns the current match ID
func (s *AppState) GetMatchID() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.matchID
}

// SetMatchID sets the current match ID
func (s *AppState) SetMatchID(matchID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.matchID = matchID
}

// IsAuthenticated checks if the user is authenticated
func (s *AppState) IsAuthenticated() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.userID != "" && s.username != ""
}

// Reset clears all state values
func (s *AppState) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.userID = ""
	s.username = ""
	s.matchID = ""
}