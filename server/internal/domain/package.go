package domain

import (
	"cod-server/internal/utils"
	"errors"
)

type Package struct {
	ID	string				`json:"id"`
	Cards	utils.Map[string, *Card]	`json:"cards"`
}

func (pkg *Package) Unpack() (*Card, *Card, *Card, error) {
	var ok bool
	var rock, paper, scissors *Card
	rock, ok = pkg.Cards.Get("rock")
	if !ok {
		return nil, nil, nil, errors.New("rock card not found in package")
	}
	paper, ok = pkg.Cards.Get("paper")
	if !ok {
		return nil, nil, nil, errors.New("paper card not found in package")
	}
	scissors, ok = pkg.Cards.Get("scissors")
	if !ok {
		return nil, nil, nil, errors.New("scissors card not found in package")
	}
	return rock, paper, scissors, nil
}
