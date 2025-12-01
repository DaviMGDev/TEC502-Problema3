package services

import (
	"cod-server/internal/domain"
	"cod-server/internal/data"
)

type RoomServiceInterface interface {
	CreateRoom(hostID string) (domain.Room, error)
	JoinRoom(roomID, userID string) error
	LeaveRoom(roomID, userID string) error
	GetRoomInfo(roomID string) (domain.Room, error)
	KickUser(roomID, hostID, userID string) error 
	SearchForRoom(userID string) (domain.Room, error)
}

type RoomService struct {
	roomRepo data.Repository[domain.Room]
	userRepo data.Repository[domain.User]
}
