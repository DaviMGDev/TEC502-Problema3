package services

import (
	"game-server/internal/domain"
	"game-server/internal/storage"
)

type StoreService struct {
	cardsRepository storage.Repository[domain.Card]
	userRepository  storage.Repository[domain.User]
	roomRepository  storage.Repository[domain.Room]
}

func NewStoreService(
	cardsRepo storage.Repository[domain.Card],
	userRepo storage.Repository[domain.User],
	roomRepo storage.Repository[domain.Room],
) *StoreService {
	return &StoreService{
		cardsRepository: cardsRepo,
		userRepository:  userRepo,
		roomRepository:  roomRepo,
	}
}

func (s *StoreService) GetAllCards() ([]domain.Card, error) {
	return s.cardsRepository.GetAll()
}

func (s *StoreService) PurchaseCard(userID, cardID string) error {
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return err
	}

	card, err := s.cardsRepository.GetByID(cardID)
	if err != nil {
		return err
	}

	if user.Balance < card.Price {
		return domain.ErrInsufficientFunds
	}

	user.Balance -= card.Price
	user.Cards = append(user.Cards, card)
