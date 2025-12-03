package services

import (
	"cod-server/internal/domain"
)

type StoreService interface {
	ReplenishPackages() error 
	GetAvailablePackages(userID string) ([]domain.Package, error)
	BuyPackage(userID string, packageID string) error 
}

