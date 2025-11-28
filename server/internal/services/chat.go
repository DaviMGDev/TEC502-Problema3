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

