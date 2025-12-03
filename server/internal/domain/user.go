package domain

import "cod-server/internal/utils"

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	PassHash string `json:"-"`
	Deck utils.Map[string, *Card] `json:"deck"`
}
