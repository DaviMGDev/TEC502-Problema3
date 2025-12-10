package domain

import (

	"log"

	"testing"

)



func TestCard_Against(t *testing.T) {

	// Test cases for different types

	tests := []struct {

		name     string

		card     *Card

		opponent *Card

		expected int

	}{

		// Different Types (RPS Logic)

		{"Rock vs Scissors (Win)", &Card{ID: "c1", Type: Rock, Level: 1}, &Card{ID: "c2", Type: Scissors, Level: 1}, 1},

		{"Paper vs Rock (Win)", &Card{ID: "c1", Type: Paper, Level: 1}, &Card{ID: "c2", Type: Rock, Level: 1}, 1},

		{"Scissors vs Paper (Win)", &Card{ID: "c1", Type: Scissors, Level: 1}, &Card{ID: "c2", Type: Paper, Level: 1}, 1},



		{"Rock vs Paper (Lose)", &Card{ID: "c1", Type: Rock, Level: 1}, &Card{ID: "c2", Type: Paper, Level: 1}, -1},

		{"Paper vs Scissors (Lose)", &Card{ID: "c1", Type: Paper, Level: 1}, &Card{ID: "c2", Type: Scissors, Level: 1}, -1},

		{"Scissors vs Rock (Lose)", &Card{ID: "c1", Type: Scissors, Level: 1}, &Card{ID: "c2", Type: Rock, Level: 1}, -1},



		// Same Types (Level Tie-breaking)

		{"Rock vs Rock (Same Level - Draw)", &Card{ID: "c1", Type: Rock, Level: 5}, &Card{ID: "c2", Type: Rock, Level: 5}, 0},

		{"Rock vs Rock (Higher Level Wins)", &Card{ID: "c1", Type: Rock, Level: 10}, &Card{ID: "c2", Type: Rock, Level: 5}, 1},

		{"Rock vs Rock (Lower Level Loses)", &Card{ID: "c1", Type: Rock, Level: 5}, &Card{ID: "c2", Type: Rock, Level: 10}, -1},



		{"Paper vs Paper (Same Level - Draw)", &Card{ID: "c1", Type: Paper, Level: 5}, &Card{ID: "c2", Type: Paper, Level: 5}, 0},

		{"Paper vs Paper (Higher Level Wins)", &Card{ID: "c1", Type: Paper, Level: 10}, &Card{ID: "c2", Type: Paper, Level: 5}, 1},

		{"Paper vs Paper (Lower Level Loses)", &Card{ID: "c1", Type: Paper, Level: 5}, &Card{ID: "c2", Type: Paper, Level: 10}, -1},



		{"Scissors vs Scissors (Same Level - Draw)", &Card{ID: "c1", Type: Scissors, Level: 5}, &Card{ID: "c2", Type: Scissors, Level: 5}, 0},

		{"Scissors vs Scissors (Higher Level Wins)", &Card{ID: "c1", Type: Scissors, Level: 10}, &Card{ID: "c2", Type: Scissors, Level: 5}, 1},

		{"Scissors vs Scissors (Lower Level Loses)", &Card{ID: "c1", Type: Scissors, Level: 5}, &Card{ID: "c2", Type: Scissors, Level: 10}, -1},

	}



	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			log.Printf("Running test: %s", tt.name)

			result := tt.card.Against(tt.opponent)

			if result != tt.expected {

				t.Errorf("card: %+v, opponent: %+v, Expected: %d, Got: %d", tt.card, tt.opponent, tt.expected, result)

			} else {

				log.Printf("Test '%s' passed.", tt.name)

			}

		})

	}

}


