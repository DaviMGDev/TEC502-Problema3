package domain

const (
	Rock		= "rock"
	Paper		= "paper"
	Scissors	= "scissors"
)

type Card struct {
	ID	string	`json:"id"`
	Type	string	`json:"type"`
	Level	int64	`json:"level"`
}

func (card *Card) Against(opponent *Card) int {
	if card.Type == opponent.Type {

		if card.Level > opponent.Level {
			return 1
		} else if card.Level < opponent.Level {
			return -1
		}
		return 0
	}
	if (card.Type == Rock && opponent.Type == Scissors) ||
		(card.Type == Paper && opponent.Type == Rock) ||
		(card.Type == Scissors && opponent.Type == Paper) {
		return 1
	}
	return -1
}
