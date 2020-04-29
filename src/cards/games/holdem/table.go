package holdem

import (
	"fmt"
	"qpoker/utils"
)

// Table holds the information about a given card table
type Table struct {
	Players  []*Player `json:"players"`
	Capacity int       `json:"capacity"`

	// Used for internal hand logic
	activeIndex int
	dealerIndex int

	// Public facing for clients
	Active int64 `json:"active"`
	Dealer int64 `json:"dealer"`
}

// NewTable returns a new table object
func NewTable(capacity int, players []*Player) *Table {
	tablePlayers := make([]*Player, capacity)

	for i := range players {
		if i >= capacity {
			break
		}

		tablePlayers[i] = players[i]
	}

	return &Table{
		Players:     tablePlayers,
		Capacity:    capacity,
		activeIndex: 0,
		dealerIndex: 0, // TODO: Random?
	}
}

func (t *Table) next(pos int) int {
	return (pos + 1) % t.Capacity
}

func (t *Table) nextPos(pos int) int {
	next := pos
	for {
		next = t.next(next)
		if next == pos {
			return -1
		}

		if next == pos {
			fmt.Printf("Error: next pos could not find a next: %d, %d\n", pos, t.Capacity)
			return pos // TODO: this is an error state
		}

		if t.Players[next] == nil || !t.Players[next].IsActive() {
			continue
		}

		break
	}

	return next
}

func (t *Table) getPlayerID(pos int) int64 {
	return t.Players[pos].ID
}

// GetActivePlayer returns the currently active player
func (t *Table) GetActivePlayer() *Player {
	return t.Players[t.activeIndex]
}

// GetActivePlayers returns the currently active players ordered by turn
func (t *Table) GetActivePlayers() []*Player {
	players := []*Player{}
	playerIndexes := []int{}
	start := t.activeIndex
	next := t.nextPos(start)

	for next != -1 && !utils.IntSliceHasValue(playerIndexes, next) {
		players = append(players, t.Players[next])
		playerIndexes = append(playerIndexes, next)
		next = t.nextPos(next)
	}

	return players
}

// ActivateNextPlayer moves the table focus to the next player
func (t *Table) ActivateNextPlayer(getActions func() map[string]bool) {
	t.GetActivePlayer().SetPlayerActions(nil)
	t.activeIndex = t.nextPos(t.activeIndex)
	t.Active = t.getPlayerID(t.activeIndex)
	t.GetActivePlayer().SetPlayerActions(getActions())
}

// NextRound moves the table focus to first bet for a new betting round
func (t *Table) NextRound(getActions func() map[string]bool) {
	t.resetPlayerRoundStates()
	t.activeIndex = t.dealerIndex
	t.ActivateNextPlayer(getActions)
}

// NextHand prepares table for next hand
func (t *Table) NextHand() error {
	t.resetPlayerHandStates()

	// Boot invalid stacks to pending
	for i := range t.Players {
		if t.Players[i] != nil && t.Players[i].Stack == int64(0) {
			t.Players[i].State = PlayerStatePending
		}
	}

	activePlayers := t.GetActivePlayers()
	if len(activePlayers) < 2 {
		return fmt.Errorf("Not enough active players: %+v", activePlayers)
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

// GetAllPlayers retrieves all players at table
func (t *Table) GetAllPlayers() []*Player {
	players := []*Player{}

	for i := range t.Players {
		if t.Players[i] == nil {
			continue
		}

		if t.Players[i].State == PlayerStateInit && t.Players[i].Stack <= int64(0) {
			continue
		}

		players = append(players, t.Players[i])
	}

	return players
}

func (t *Table) resetPlayerRoundStates() {
	for i := range t.Players {
		if t.Players[i] == nil {
			continue
		}

		if !t.Players[i].IsActive() {
			continue
		}

		t.Players[i].State = PlayerStateInit
		t.Players[i].SetPlayerActions(nil)
	}
}

func (t *Table) resetPlayerHandStates() {
	for i := range t.Players {
		if t.Players[i] == nil {
			continue
		}

		t.Players[i].State = PlayerStateInit
		t.Players[i].SetPlayerActions(nil)
		t.Players[i].BigBlind = false
		t.Players[i].LittleBlind = false
	}
}
