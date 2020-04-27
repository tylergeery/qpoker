package holdem

import "qpoker/cards"

const (
	// PlayerStateInit is the init state for each hand
	PlayerStateInit = "init"
	// PlayerStateFold means a player has folded
	PlayerStateFold = "fold"
	// PlayerStateCall means a player has called
	PlayerStateCall = "call"
	// PlayerStatePending means a player is pending for next round
	PlayerStatePending = "pending"
	// PlayerStateCheck means a player has checked
	PlayerStateCheck = "check"
	// PlayerStateBet means a player has bet
	PlayerStateBet = "bet"
)

// Player holds the information about a player at a table
type Player struct {
	ID           int64           `json:"id"`
	Username     string          `json:"username"`
	Cards        []cards.Card    `json:"-"`
	CardsVisible bool            `json:"-"`
	Stack        int64           `json:"stack"`
	Options      map[string]bool `json:"options"`
	State        string          `json:"state"`
	BigBlind     bool            `json:"big_blind"`
	LittleBlind  bool            `json:"little_blind"`
}

// SetPlayerActions sets the moves a player is allowed to make
func (p *Player) SetPlayerActions(actions map[string]bool) {
	p.Options = actions
}

// IsActive returns whether the player is active in the current hand
func (p *Player) IsActive() bool {
	if p.State == PlayerStatePending || p.State == PlayerStateFold {
		return false
	}

	return true
}
