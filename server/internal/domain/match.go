package domain

type Match struct {
	ID	string		`json:"id"`
	Players	[2]string	`json:"players"`
	Moves	[2]*Card	`json:"moves"`
	Winner	string		`json:"winner"`
}
