package domain

import "cod-server/internal/utils"

type User struct {
	ID 		 string `json:"id"` 
	Username string `json:"username"` 
	Password string `json:"password"` 
	Cards utils.Map[string, *Card] `json:"cards"`
}
