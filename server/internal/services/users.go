package services

import (
	"cod-server/internal/domain"
	"cod-server/internal/data"
	"sync"
)

var (
	userCounter uint 
	userCounterMutex sync.RWMutex
)

func init() {
	userCounter = 0
	userCounterMutex = sync.RWMutex{}
}

type UserServiceInterface interface {
	Register(username, password string) (error)
	Login(username, password string) (string, error)
	GetUserName(userID int) (string, error)
	ChangePassword(userID int, oldPassword, newPassword string) (error)
}

type UserService struct {
	userRepo data.Repository[domain.User]
}
