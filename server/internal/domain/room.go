package domain

type Room struct {
	ID      string   `json:"id"` 
	Members []string `json:"members"`
}
