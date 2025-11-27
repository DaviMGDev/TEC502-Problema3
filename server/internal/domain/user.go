package domain

type User struct {
	ID 						 string `json:"id"`
	Name 					 string `json:"name"`
	Username 			 string `json:"username"`
	HashedPassword string `json:"-"`
}
