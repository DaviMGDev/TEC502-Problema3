package services

import (
	"game-server/internal/domain"
	"game-server/internal/storage"
)

type UserService struct {
	repo storage.Repository[domain.User] 
}

func NewUserService(repo storage.Repository[domain.User]) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(user domain.User) error {
	return s.repo.Create(user.ID, user)
}

func (s *UserService) GetUser(id string) (domain.User, error) {
	return s.repo.Read(id)
}

func (s *UserService) UpdateUser(user domain.User) error {
	return s.repo.Update(user.ID, user)
}

func (s *UserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

func (s *UserService) ListUsers() ([]domain.User, error) {
	return s.repo.List()
}

