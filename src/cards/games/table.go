package games

import (
	"fmt"
)

// Table holds the information about a given card table
type Table struct {
	Players  []*Player `json:"players"`
	Capacity int       `json:"capacity"`
}

// AddPlayer to table
func (t *Table) AddPlayer(player *Player) error {
	if len(t.Players) == t.Capacity {
		return fmt.Errorf("Table is full at %d players", t.Capacity)
	}

	t.Players = append(t.Players, player)

	return nil
}

// RemovePlayer from table
func (t *Table) RemovePlayer(playerID int64) {
	for i := range t.Players {
		if t.Players[i].ID == playerID {
			t.Players = append(t.Players[:i], t.Players[i+1:]...)
			return
		}
	}
}
