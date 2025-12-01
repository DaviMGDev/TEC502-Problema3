package domain

type CardType uint8 

const (
	Rock CardType = iota 
	Paper 
	Scissors
)

type Card struct {
	ID   string   `json:"id"`
	Type CardType `json:"type"`
	Level uint8    `json:"level"`
}
