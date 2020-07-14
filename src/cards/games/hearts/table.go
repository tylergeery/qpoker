package holdem

import (
	"qpoker/cards"
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
	heartsTable := &Table{
		table,
	}

	for i := range players {
		heartsTable.AddPlayer(players[i])
	}

	return heartsTable
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

// NextHand advances to the next hand
func (t *Table) NextHand() {
	// reset players cards
	players := t.GetAllPlayers()
	for i := range players {
		players[i].Cards = []cards.Card{}
		players[i].Pile = []cards.Card{}
	}
	t.DealerIndex = t.NextPos(t.DealerIndex)
	t.Dealer = t.GetPlayerID(t.DealerIndex)
	t.ActiveIndex = t.NextPos(t.DealerIndex)
	t.Active = t.GetPlayerID(t.ActiveIndex)
}
