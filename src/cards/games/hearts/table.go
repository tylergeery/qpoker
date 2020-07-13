package holdem

import (
	"fmt"
	"time"
)

// Table holds the information about a given card table
type Table struct {
	Players  []*Player `json:"players"`
	Capacity int       `json:"capacity"`

	// Used for internal hand logic
	activeIndex int
	dealerIndex int

	// Public facing for clients
	Active   int64 `json:"active"`
	ActiveAt int64 `json:"active_at"`
	Dealer   int64 `json:"dealer"`
}

// NewTable returns a new table object
func NewTable(players []*Player) *Table {
	tablePlayers := make([]*Player, 4)
	for i := range players {
		tablePlayers[i] = players[i]
	}

	return &Table{
		Players:     tablePlayers,
		Capacity:    4,
		activeIndex: 0,
		dealerIndex: 0,
	}
}

func (t *Table) next(pos int) int {
	return (pos + 1) % t.Capacity
}

func (t *Table) nextPos(pos int) int {
	return t.next(pos)
}

func (t *Table) getPlayerID(pos int) int64 {
	return t.Players[pos].ID
}

// GetActivePlayer returns the currently active player
func (t *Table) GetActivePlayer() *Player {
	return t.Players[t.activeIndex]
}

// GetPlayers returns the currently active players ordered by turn
func (t *Table) GetPlayers() []*Player {
	return t.Players
}

// GetActivePlayers returns the currently active players ordered by turn
func (t *Table) GetActivePlayers() []*Player {
	players := []*Player{}

	for i := range t.Players {
		if t.Players[i] != nil {
			players = append(players, t.Players[i])
		}
	}

	return players
}

// ActivateNextPlayer moves the table focus to the next player
func (t *Table) ActivateNextPlayer(getActions func() map[string]bool) {
	t.GetActivePlayer().SetPlayerActions(nil)
	t.activeIndex = t.nextPos(t.activeIndex)
	t.Active = t.getPlayerID(t.activeIndex)
	t.ActiveAt = time.Now().Unix()
	t.GetActivePlayer().SetPlayerActions(getActions())
}

// NextHand prepares table for next hand
func (t *Table) NextHand() error {
	t.resetPlayerHandStates()

	// Boot invalid stacks to pending
	for i := range t.Players {
		if t.Players[i] == nil {
			return fmt.Errorf("Not enough active players: %+v", t.Players)
		}
	}

	t.dealerIndex = t.nextPos(t.dealerIndex)
	t.Dealer = t.getPlayerID(t.dealerIndex)
	t.activeIndex = t.nextPos(t.dealerIndex)
	t.Active = t.getPlayerID(t.activeIndex)

	return nil
}

// AddPlayer to table
func (t *Table) AddPlayer(player *Player) error {
	for i := range t.Players {
		if t.Players[i] != nil && t.Players[i].ID == player.ID {
			return fmt.Errorf("Player %d already exists at table", player.ID)
		}
	}

	for i := range t.Players {
		if t.Players[i] == nil {
			t.Players[i] = player
			return nil
		}
	}

	return fmt.Errorf("Table is full at %d players", t.Capacity)
}

// RemovePlayer from table
func (t *Table) RemovePlayer(playerID int64) error {
	for i := range t.Players {
		if t.Players[i] != nil && t.Players[i].ID == playerID {
			t.Players[i] = nil
			return nil
		}
	}

	return fmt.Errorf("Player %d was not found at the table", playerID)
}

func (t *Table) resetPlayerHandStates() {
	for i := range t.Players {
		if t.Players[i] == nil {
			continue
		}

		t.Players[i].SetPlayerActions(nil)
	}
}
