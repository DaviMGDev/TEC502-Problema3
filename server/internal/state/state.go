package state

import (
	// "cod-server/internal/data"
	// "cod-server/internal/domain"
	// "cod-server/internal/services"
)

type State struct {
	Address string 
	BrokerAddress string 
	IsLeader bool 
}

func NewState() *State {
	return &State{
		Address:      "",
		BrokerAddress: "",
		IsLeader:     false,
	}
}
