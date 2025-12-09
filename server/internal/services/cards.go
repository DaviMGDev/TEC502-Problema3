package services

import (
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"cod-server/internal/utils"
	"math/rand"
	"strconv"
)

type CardService interface {
	BuyCardPack(userID string) error 
	SwapCard(user1, user2, cardType string) error 
	ListUserCards(userID string) ([]*domain.Card, error)
}

type CardServiceImplementation struct {
	userRepo data.Repository[*domain.User]
	packages utils.List[*domain.Package]
	packCounter uint64
}

func NewCardService(userRepo data.Repository[*domain.User]) *CardServiceImplementation {
	service := &CardServiceImplementation{
		userRepo: userRepo,
		packages: utils.NewSafeList[*domain.Package](),
	}
	service.checkStock()
	return service
}

func (s *CardServiceImplementation) checkStock() {
	if s.packages.Size() < 16 {
		for i := 0; i < 64; i++ {
			pack := domain.Package{
				ID: strconv.FormatUint(s.packCounter, 16),
				Cards: utils.NewSafeMap[string, *domain.Card](),
			}
			s.packCounter++
			types := []string{"rock", "paper", "scissors"}
			level := (rand.Int63() % 100) + 1
			for _, t := range types {	
				card := &domain.Card{
					ID: strconv.FormatUint(rand.Uint64(), 16),
					Type: t,
					Level: level,
				}
				pack.Cards.Set(t, card)
			}
			s.packages.Append(&pack)
		}
	}
}

func (s *CardServiceImplementation) BuyCardPack(userID string) error {
	defer s.checkStock()
	user, err := s.userRepo.Read(userID)
	if err != nil {
		return err
	}
	pack, ok := s.packages.Pop()
	if !ok {
		return nil 
	}
	for _, card := range pack.Cards.Values() {
		user.Cards.Set(card.Type, card)
	}
	return s.userRepo.Update(userID, user)
}

func (s *CardServiceImplementation) SwapCard(user1, user2, cardType string) error {
	first, err := s.userRepo.Read(user1)
	if err != nil {
		return err 
	}
	second, err := s.userRepo.Read(user2)
	if err != nil {
		return err 
	}
	card1, ok1 := first.Cards.Get(cardType)
	card2, ok2 := second.Cards.Get(cardType)
	if !ok1 || !ok2 {
		return nil 
	}
	first.Cards.Set(cardType, card2)
	second.Cards.Set(cardType, card1)
	if err := s.userRepo.Update(user1, first); err != nil {
		return err 
	}
	if err := s.userRepo.Update(user2, second); err != nil {
		return err 
	}
	return nil
}

func (s *CardServiceImplementation) ListUserCards(userID string) ([]*domain.Card, error) {
	user, err := s.userRepo.Read(userID)
	if err != nil {
		return nil, err 
	}
	return user.Cards.Values(), nil 
}
