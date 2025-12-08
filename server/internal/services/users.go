package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"errors"
	"strconv"
)

type UserService interface {
	Register(username, password string) error 
	Login(username, password string) (string, error) 
}

type UserServiceImplementation struct {
	userRepo data.Repository[*domain.User]
	userCounter uint64
}

func NewUserService(userRepo data.Repository[*domain.User]) *UserServiceImplementation {
	return &UserServiceImplementation{
		userRepo: userRepo,
	}
}

func (s *UserServiceImplementation) Register(username, password string) error {
	user := &domain.User{
		ID: strconv.FormatUint(s.userCounter, 16),
		Username: username,
		Password: password,
	}
	s.userCounter++
	return s.userRepo.Create(username, user)
}

func (s *UserServiceImplementation) Login(username, password string) (string, error) {
	user, err := s.userRepo.Read(username)
	if err != nil {
		return "", err
	}
	if user.Password != password {
		return "", errors.New("invalid credentials")
	}
	return user.ID, nil 
}
