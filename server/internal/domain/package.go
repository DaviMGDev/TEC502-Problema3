package domain

type Package struct {
	ID    string   `json:"id"`
	Cards []*Card `json:"cards"`
}
