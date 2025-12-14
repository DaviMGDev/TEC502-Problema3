package domain 

import (
	"cod-server/internal/utils"
)

type User struct {
	ID       utils.String `json:"id"`	
	Username utils.String `json:"username"`
	Password utils.String `json:"password"`
	Hand *Pack `json:"hand"`
}

func NewUser(id, username, password utils.String, hand *Pack) *User {
	return &User{
		ID:       id,
		Username: username,
		Password: password,
		Hand:     hand,
	}
}

func (u *User) Equals(other utils.Comparable) bool {
	otherUser, ok := other.(*User)
	if !ok {
		return false
	}
	return u.ID == otherUser.ID &&
		u.Username == otherUser.Username &&
		u.Password == otherUser.Password &&
		u.Hand.Equals(otherUser.Hand)	
}
