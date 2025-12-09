package domain 

const (
	Rock = "rock"
	Paper = "paper"
	Scissors = "scissors"
)

type Card struct {
	ID  string   `json:"id"` 
	Type string `json:"type"` 
	Level int64 	`json:"level"`
}

func (card *Card) Against(opponent *Card) int {
	if card.Type == opponent.Type {
		// Tie-breaking by Level
		if card.Level > opponent.Level {
			return 1 // card wins by level
		} else if card.Level < opponent.Level {
			return -1 // card loses by level
		}
		return 0 // Draw (same type, same level)
	}
	if (card.Type == Rock && opponent.Type == Scissors) ||
		(card.Type == Paper && opponent.Type == Rock) ||
		(card.Type == Scissors && opponent.Type == Paper) {
		return 1 
	}
	return -1 
}
