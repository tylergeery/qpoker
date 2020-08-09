package holdem

import (
	"fmt"
	"qpoker/cards/games"
)

// Table holds the information about a given card table
type Table struct {
	*games.Table
}

// NewTable returns a new table object
func NewTable(capacity int, players []*Player) *Table {
	if capacity < len(players) {
		capacity = len(players)
	}

	table := games.NewTable(capacity)
	holdemTable := &Table{
		table,
	}

	for i := range players {
		holdemTable.AddPlayer(players[i])
	}

	return holdemTable
}

// GetActivePlayer returns the currently active player
func (t *Table) GetActivePlayer() *Player {
	return t.GetActiveIPlayer().(*Player)
}

// GetPlayerByID returns player with id
func (t *Table) GetPlayerByID(id int64) *Player {
	return t.GetIPlayerByID(id).(*Player)
}

func getPlayers(iPlayers []games.IPlayer) []*Player {
	players := make([]*Player, len(iPlayers))
	for i := range iPlayers {
		players[i] = iPlayers[i].(*Player)
	}

	return players
}

// GetAllPlayers returns the currently active players ordered by turn
func (t *Table) GetAllPlayers() []*Player {
	return getPlayers(t.GetAllIPlayers())
}

// GetActivePlayers returns the currently active players ordered by turn
func (t *Table) GetActivePlayers() []*Player {
	return getPlayers(t.GetActiveIPlayers())
}

// GetReadyPlayers return the players ready to play
func (t *Table) GetReadyPlayers() []*Player {
	return getPlayers(t.GetReadyIPlayers())
}

// GetPlayersFromIDs return all IPlayers with ID
func (t *Table) GetPlayersFromIDs(ids []int64) []*Player {
	return getPlayers(t.GetIPlayersFromIDs(ids))
}

// NextRound moves the table focus to first bet for a new betting round
func (t *Table) NextRound(getActions func() map[string]bool) {
	t.resetPlayerRoundStates()
	t.ActiveIndex = t.DealerIndex
	t.ActivateNextPlayer(getActions)
}

// NextHand prepares table for next hand
func (t *Table) NextHand() error {
	t.resetPlayerHandStates()

	// Boot invalid stacks to pending
	for _, player := range t.GetAllPlayers() {
		if player.Stack == int64(0) {
			player.State = PlayerStatePending
		}
	}

	activePlayers := t.GetActivePlayers()
	if len(activePlayers) < 2 {
		return fmt.Errorf("Not enough active players: %+v", activePlayers)
	}

	t.DealerIndex = t.NextPos(t.DealerIndex)
	t.Dealer = t.GetPlayerID(t.DealerIndex)
	t.ActiveIndex = t.NextPos(t.DealerIndex)
	t.Active = t.GetPlayerID(t.ActiveIndex)

	return nil
}

func (t *Table) resetPlayerRoundStates() {
	for _, player := range t.GetAllPlayers() {
		if !player.IsActive() {
			continue
		}

		player.State = PlayerStateInit
		player.SetPlayerActions(nil)
	}
}

func (t *Table) resetPlayerHandStates() {
	for _, player := range t.GetAllPlayers() {
		player.State = PlayerStateInit
		player.SetPlayerActions(nil)
		player.BigBlind = false
		player.LittleBlind = false
	}
}
