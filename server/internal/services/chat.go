package services

import (
	"game-server/internal/domain"
	"game-server/internal/storage"
)

type ChatService struct {
	messageRepo storage.Repository[domain.Message]
	roomRepo storage.Repository[domain.Room]
	userRepo storage.Repository[domain.User]
}

func NewChatService(
	messageRepository storage.Repository[domain.Message],
	roomRepository storage.Repository[domain.Room],
	userRepository storage.Repository[domain.User],
) *ChatService {
	return &ChatService{
		messageRepo: messageRepository,
		roomRepo: roomRepository,
		userRepo: userRepository,
	}
}

func (s *ChatService) SendMessage(roomID string, message domain.Message) error {
	// Check if room exists
	_, err := s.roomRepo.Read(roomID)
	if err != nil {
		return err
	}

	// Store the message
	return s.messageRepo.Create(message.Timestamp.String(), message)
}

func (s *ChatService) GetMessages() ([]domain.Message, error) {
	return s.messageRepo.List()
}

func (s *ChatService) ListRooms() ([]domain.Room, error) {
	return s.roomRepo.List()
}

func (s *ChatService) ListUsers() ([]domain.User, error) {
	return s.userRepo.List()
}

