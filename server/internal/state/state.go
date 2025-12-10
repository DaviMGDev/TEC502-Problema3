package state

import ()

type State struct {
	Address		string
	BrokerAddress	string
}

func NewState() *State {
	return &State{
		Address:	"",
		BrokerAddress:	"",
	}
}
