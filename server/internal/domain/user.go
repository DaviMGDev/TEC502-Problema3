package domain

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	PassHash string `json:"-"`
	Deck *Package `json:"deck"`
}
