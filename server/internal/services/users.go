package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"errors"
	"strconv"
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

type UserService interface {
	Register(username, password string) (error)
	Login(username, password string) (string, error)
}

