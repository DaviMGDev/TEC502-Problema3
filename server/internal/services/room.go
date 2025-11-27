package services

import (
	"game-server/internal/domain"
	"game-server/internal/storage"
)

type RoomService struct {
	repo storage.Repository[domain.Room]
}

func NewRoomService(repository storage.Repository[domain.Room]) *RoomService {
	return &RoomService{
		repo: repository,
	}
}

func (s *RoomService) CreateRoom(room domain.Room) error {
	return s.repo.Create(room.ID, room)
}

func (s *RoomService) GetRoom(id string) (domain.Room, error) {
	return s.repo.Read(id)
}

func (s *RoomService) UpdateRoom(room domain.Room) error {
	return s.repo.Update(room.ID, room)
}

func (s *RoomService) DeleteRoom(id string) error {
	return s.repo.Delete(id)
}

func (s *RoomService) ListRooms() ([]domain.Room, error) {
	return s.repo.List()
}

func (s *RoomService) AddPlayerToRoom(roomID string, playerID string) error {
	room, err := s.repo.Read(roomID)
	if err != nil {
		return err
	}

	if room.Count >= 2{
		return nil // Room is full
	}

	room.Users = append(room.Users, playerID)
	room.Count++

	return s.repo.Update(roomID, room)
}

func (s *RoomService) RemovePlayerFromRoom(roomID string, playerID string) error {
	room, err := s.repo.Read(roomID)
	if err != nil {
		return err
	}

	for i, id := range room.Users {
		if id == playerID {
			room.Users = append(room.Users[:i], room.Users[i+1:]...)
			room.Count--
			break
		}
	}

	return s.repo.Update(roomID, room)
}

func (s *RoomService) IsRoomFull(roomID string) (bool, error) {
	room, err := s.repo.Read(roomID)
	if err != nil {
		return false, err
	}

	return room.Count >= 2, nil
}

