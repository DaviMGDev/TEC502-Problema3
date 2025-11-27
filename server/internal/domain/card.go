package domain

const (
	ROCK uint8 = iota 
	PAPER
	SCISSORS
)

type Card struct {
	ID 		string `json:"id"`
	Type 	uint8 `json:"type"`
	Level uint8 `json:"level"`
}
