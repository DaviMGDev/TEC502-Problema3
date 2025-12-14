package services

import "cod-server/internal/domain"

type UserServiceInterface interface {
	Register(username, password string) error 
	Login(username, password string) (*domain.User, error)
}
