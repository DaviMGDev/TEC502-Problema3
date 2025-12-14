package domain

import "cod-server/internal/utils"

type Card struct {
	ID   utils.String `json:"id"`
	Type utils.String `json:"type"`
}

func NewCard(id utils.String, cardType utils.String) *Card {
	return &Card{
		ID: id,
		Type: cardType,
	}
}

func (c *Card) Equals(other utils.Comparable) bool {
	otherCard, ok := other.(*Card)
	if !ok {
		return false
	}
	return c.ID == otherCard.ID && c.Type == otherCard.Type
}

func (c *Card) Against(opponentCard *Card) int {
	if c.Type == opponentCard.Type {
		return 0 // Draw
	}

	if (c.Type == "rock" && opponentCard.Type == "scissors") ||
		(c.Type == "scissors" && opponentCard.Type == "paper") ||
		(c.Type == "paper" && opponentCard.Type == "rock") {
		return 1 // Win
	}

	return -1 // Lose
}

type Pack struct {
	ID utils.String `json:"id"` 
	Cards utils.List[*Card] `json:"cards"`
}

func (p *Pack) Equals(other utils.Comparable) bool {
	otherPack, ok := other.(*Pack)
	if !ok {
		return false
	}
	if p.ID != otherPack.ID || p.Cards.Size() != otherPack.Cards.Size() {
		return false
	}
	for i := 0; i < p.Cards.Size(); i++ {
		card1, _ := p.Cards.Get(i)
		card2, _ := otherPack.Cards.Get(i)
		if !card1.Equals(card2) {
			return false
		}
	}
	return true
}

func NewPack(id utils.String, cards utils.List[*Card]) *Pack {
	return &Pack{
		ID:    id,
		Cards: cards,
	}
}

func (p *Pack) DrawCard(index int) (*Card, error) {
	return p.Cards.Extract(index)
}

func (p *Pack) AddCard(card *Card) {
	p.Cards.Append(card)
}


