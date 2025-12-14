package domain

import "cod-server/internal/utils"

type Match struct {
	ID         utils.String   `json:"id"`          // Match ID
	Players    utils.List[utils.String] `json:"players"`     // struct.User.ID
	Clients    utils.List[utils.String] `json:"clients"`     // ClientIDs
	Moves      utils.List[utils.Map[utils.String, *Card]] `json:"moves"`       // [{struct.User.ID: *Card}]
	Winner     utils.String   `json:"winner"`      // struct.User.ID
}

func NewMatch(id utils.String, players utils.List[utils.String], clients utils.List[utils.String]) *Match {
	match := &Match{
		ID:      id,
		Players: players,
		Clients: clients,
		Winner:  "",
	}	
	match.Moves = utils.NewSafeList[utils.Map[utils.String, *Card]]()

	return match
}

func (m *Match) Equals(other utils.Comparable) bool {
	otherMatch, ok := other.(*Match)
	if !ok {
		return false
	}
	return m.ID.Equals(otherMatch.ID) && m.Players.Equals(otherMatch.Players) && m.Clients.Equals(otherMatch.Clients) && m.Moves.Equals(otherMatch.Moves) && m.Winner.Equals(otherMatch.Winner)
}

func (m *Match) AddPlayer(playerID utils.String, clientID utils.String) {
	if m.Players.Size() >= 2 {
		return
	}
	m.Players.Append(playerID)
	m.Clients.Append(clientID)
}

func (m *Match) RemovePlayer(playerID utils.String) {
	index := m.Players.Search(playerID)
	if index == -1 {
		return
	}
	m.Players.Extract(index)
	m.Clients.Extract(index)
}

func (m *Match) MakeMove(playerID utils.String, card *Card) {
	previusMove := m.Moves.Filter(func(move utils.Map[utils.String, *Card]) bool {
		_, exists := move.Get(playerID)
		return exists
	})
	if previusMove.Size() > 0 {
		return
	}
	if m.Moves.Size() >= 2 {
		return
	}
	if m.Players.Search(playerID) == -1 {
		return
	}
	move := utils.NewSafeMap[utils.String, *Card]()
	move.Set(playerID, card)
	m.Moves.Append(move)
}
