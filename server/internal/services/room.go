package services

import (
)

var (
//	roomCounter uint
)

type RoomService interface {
	CreateRoom(hostID string) (string, error)
	JoinRoom(roomID, userID string) error
	LeaveRoom(roomID, userID string) error
	KickUser(roomID, hostID, userID string) error 
	SearchForRoom(userID string) (string, error)
}

