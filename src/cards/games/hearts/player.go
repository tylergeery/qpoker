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

// GetID return player ID
func (p *Player) GetID() int64 {
	return p.ID
}

// SetPlayerActions sets the moves a player is allowed to make
func (p *Player) SetPlayerActions(actions map[string]bool) {
	p.Options = actions
}

// IsActive returns whether the player is active in the current hand
func (p *Player) IsActive() bool {
	return true
}

// IsReady returns whether the player is ready for next game
func (p *Player) IsReady() bool {
	return true
}
