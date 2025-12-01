package domain

type Match struct {
	RoomID string    `json:"room_id"`
	Plays  [][2]map[string]*Card `json:"plays"`
	Winner string    `json:"winner"`
}
