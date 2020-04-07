package games

import "qpoker/cards"

// Player holds the information about a player at a table
type Player struct {
	ID      int64           `json:"id"`
	Cards   []cards.Card    `json:"-"`
	Stack   int64           `json:"stack"`
	Options map[string]bool `json:"options"`
	Active  bool            `json:"active"`
}
