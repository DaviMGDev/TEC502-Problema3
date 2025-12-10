package state

import "testing"

func TestNewState(t *testing.T) {
	state := NewState()

	if state == nil {
		t.Fatal("Expected NewState to return a non-nil State object, got nil")
	}

	if state.Address != "" {
		t.Errorf("Expected Address to be empty, got %s", state.Address)
	}

	if state.BrokerAddress != "" {
		t.Errorf("Expected BrokerAddress to be empty, got %s", state.BrokerAddress)
	}
}
