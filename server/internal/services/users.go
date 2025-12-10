package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"cod-server/internal/utils"
	"errors"
	"strconv"
)

type UserService interface {
	Register(username, password string) (*domain.User, error)
	Login(username, password string) (string, error)
}

type UserServiceImplementation struct {
	userRepo	data.Repository[*domain.User]
	userCounter	uint64
}

func NewUserService(userRepo data.Repository[*domain.User]) *UserServiceImplementation {
	return &UserServiceImplementation{
		userRepo: userRepo,
	}
}

func (s *UserServiceImplementation) Register(username, password string) (*domain.User, error) {

	id := strconv.FormatUint(s.userCounter, 10)
	user := &domain.User{
		ID:		id,
		Username:	username,
		Password:	password,
		Cards:		utils.NewSafeMap[string, *domain.Card](),
	}

	err := s.userRepo.Create(username, user)
	if err != nil {
		return nil, err
	}

	s.userCounter++
	return user, nil
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
