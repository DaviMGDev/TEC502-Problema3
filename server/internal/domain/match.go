package domain

import "cod-server/internal/utils"

type Match struct {
	RoomID string    `json:"room_id"`
	Plays  utils.List[[2]Play] `json:"plays"`
	Winner string    `json:"winner"`
}

type Play struct {
	UserID string `json:"user_id"`
	Card   *Card  `json:"card"`
}
