package domain

type Room struct {
	ID 	  string   `json:"id"`
	Users []string `json:"users"`
	Count uint8    `json:"count"`
}
