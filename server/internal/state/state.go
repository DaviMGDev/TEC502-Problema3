package state

import (
	"cod-server/internal/domain"
	"cod-server/internal/data"
	"cod-server/internal/services"
)

var (
	UserRepository 		data.Repository[*domain.User]
	PackageRepository data.Repository[*domain.Package]
	RoomRepository 		data.Repository[*domain.Room]

	UserService  services.UserService
	StoreService services.StoreService
	RoomService  services.RoomService 
	GameService  services.GameService 
	CardService  services.CardService
)

func init() {
}


