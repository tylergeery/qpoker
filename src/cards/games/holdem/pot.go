package holdem

import (
	"qpoker/cards/games"
)

type Pot struct {
	Bets   map[int64]int64
	Totals map[int64]int64
}

// NewPot returns a new Pot instance for the current group of players
func NewPot(players []*games.Players) Pot {
	pot := Pot{
		Bets:   map[int64]int64{},
		Totals: map[int64]int64{},
	}

	for i := range players {
		pot.Bets[players[i].ID] = 0
		pot.Totals[players[i].ID] = 0
	}

	return pot
}

// AddBet records the bet for a player
func (p Pot) AddBet(playerID, amount int64) {
	p.Bets[playerID] += amount
}

// ClearBets clears all the current to the pot
func (p Pot) ClearBets() {
	for playerID, bet := range p.Bets {
		p.Totals[playerID] += bet
		p.Bets[playerID] = 0
	}
}
