package domain

import "cod-server/internal/utils"

type Room struct {
	ID      string   `json:"id"` 
	Members  utils.List[string] `json:"members"`
}
