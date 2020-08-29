package games

import (
	"fmt"
	"qpoker/qutils"
	"time"
)

// IPlayer player interface for all games
type IPlayer interface {
	GetID() int64
	IsActive() bool
	IsReady() bool
	SetPlayerActions(actions map[string]bool)
}

// Table holds the information about a given card table
type Table struct {
	Players  []IPlayer `json:"players"`
	Capacity int       `json:"capacity"`

	// Used for internal hand logic
	ActiveIndex int `json:"-"`
	DealerIndex int `json:"-"`

	// Public facing for clients
	Active   int64 `json:"active"`
	ActiveAt int64 `json:"active_at"`
	Dealer   int64 `json:"dealer"`
}

// NewTable returns a new table object
func NewTable(capacity int) *Table {
	tablePlayers := make([]IPlayer, capacity)

	return &Table{
		Players:     tablePlayers,
		Capacity:    capacity,
		ActiveIndex: 0,
		DealerIndex: 0, // TODO: Random?
	}
}

func (t *Table) next(pos int) int {
	return (pos + 1) % t.Capacity
}

// NextPos return index for position after pos at table
func (t *Table) NextPos(pos int) int {
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

// GetPlayerID returns player ID for player at pos
func (t *Table) GetPlayerID(pos int) int64 {
	return t.Players[pos].GetID()
}

// GetActiveIPlayer returns the currently active player
func (t *Table) GetActiveIPlayer() IPlayer {
	return t.Players[t.ActiveIndex]
}

// GetIPlayerByID gets player by ID
func (t *Table) GetIPlayerByID(id int64) IPlayer {
	for i := range t.Players {
		if t.Players[i] == nil {
			continue
		}

		if id == t.Players[i].GetID() {
			return t.Players[i]
		}
	}

	return nil
}

// GetAllIPlayers returns the currently active players ordered by turn
func (t *Table) GetAllIPlayers() []IPlayer {
	players := []IPlayer{}

	for i := range t.Players {
		if t.Players[i] == nil {
			continue
		}

		players = append(players, t.Players[i])
	}

	return players
}

// GetActiveIPlayers returns the currently active players ordered by turn
func (t *Table) GetActiveIPlayers() []IPlayer {
	players := []IPlayer{}
	playerIndexes := []int{}
	start := t.ActiveIndex
	next := t.NextPos(start)

	for next != -1 && !qutils.IntSliceHasValue(playerIndexes, next) {
		players = append(players, t.Players[next])
		playerIndexes = append(playerIndexes, next)
		next = t.NextPos(next)
	}

	return players
}

// GetReadyIPlayers retrieves all players at table
func (t *Table) GetReadyIPlayers() []IPlayer {
	players := []IPlayer{}

	for i := range t.Players {
		if t.Players[i] == nil {
			continue
		}

		if !t.Players[i].IsReady() {
			continue
		}

		players = append(players, t.Players[i])
	}

	return players
}

// GetIPlayersFromIDs return all IPlayers with ID
func (t *Table) GetIPlayersFromIDs(ids []int64) []IPlayer {
	players := []IPlayer{}
	allPlayers := t.GetAllIPlayers()

	for i := range allPlayers {
		if !qutils.Int64SliceHasValue(ids, allPlayers[i].GetID()) {
			continue
		}

		players = append(players, allPlayers[i])
	}

	return players
}

// ActivateNextPlayer moves the table focus to the next player
func (t *Table) ActivateNextPlayer(getActions func() map[string]bool) {
	t.GetActiveIPlayer().SetPlayerActions(nil)
	t.ActiveIndex = t.NextPos(t.ActiveIndex)
	t.Active = t.GetPlayerID(t.ActiveIndex)
	t.ActiveAt = time.Now().Unix()
	t.GetActiveIPlayer().SetPlayerActions(getActions())
}

// AddPlayer to table
func (t *Table) AddPlayer(player IPlayer) error {
	for i := range t.Players {
		if t.Players[i] != nil && t.Players[i].GetID() == player.GetID() {
			return fmt.Errorf("Player %d already exists at table", player.GetID())
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
		if t.Players[i] != nil && t.Players[i].GetID() == playerID {
			t.Players[i] = nil
			return nil
		}
	}

	return fmt.Errorf("Player %d was not found at the table", playerID)
}
