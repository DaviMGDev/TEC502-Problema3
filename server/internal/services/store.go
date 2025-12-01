package services

import (
	"cod-server/internal/domain"
	"cod-server/internal/data"
)

type StoreServiceInterface interface {
	ReplenishPackages() error 
	GetAvailablePackages(userID string) ([]domain.Package, error)
	BuyPackage(userID string, packageID string) error 
}

type StoreService struct {
	packageRepo data.Repository[domain.Package]
	userRepo    data.Repository[domain.User]
}
