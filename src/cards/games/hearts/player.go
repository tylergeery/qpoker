package holdem

import (
	"qpoker/cards"
	"qpoker/models"
)

// Player holds the information about a player at a table
type Player struct {
	ID       int64           `json:"id"`
	Username string          `json:"username"`
	Cards    []cards.Card    `json:"-"`
	Pile     []cards.Card    `json:"pile"`
	Options  map[string]bool `json:"options"`
	Score    int
}

// NewPlayer creates a new Player
func NewPlayer(player *models.Player) *Player {
	return &Player{
		ID:       player.ID,
		Username: player.Username,
	}
}

// SetPlayerActions sets the moves a player is allowed to make
func (p *Player) SetPlayerActions(actions map[string]bool) {
	p.Options = actions
}
