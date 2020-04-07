package holdem

import (
	"qpoker/cards/games"
	"qpoker/utils"
)

// Pot controls betting pot
type Pot struct {
	PlayerBets   map[int64]int64 `json:"player_bets"`
	PlayerTotals map[int64]int64 `json:"player_totals"`
	Total        int64           `json:"total"`
}

// NewPot returns a new Pot instance for the current group of players
func NewPot(players []*games.Player) *Pot {
	pot := &Pot{
		PlayerBets:   map[int64]int64{},
		PlayerTotals: map[int64]int64{},
		Total:        0,
	}

	for i := range players {
		pot.PlayerBets[players[i].ID] = 0
		pot.PlayerTotals[players[i].ID] = 0
	}

	return pot
}

// AddBet records the bet for a player
func (p *Pot) AddBet(playerID, amount int64) {
	p.PlayerBets[playerID] += amount
}

// ClearBets clears all the current to the pot
func (p *Pot) ClearBets() {
	for playerID := range p.PlayerBets {
		p.Total += p.PlayerBets[playerID]
		p.PlayerTotals[playerID] += p.PlayerBets[playerID]
		p.PlayerBets[playerID] = 0
	}
}

// GetPayouts returns the winning amounts for each user by winning priority
func (p *Pot) GetPayouts(orderedPlayers []int64) map[int64]int64 {
	p.ClearBets()

	payouts := map[int64]int64{}
	remaining := p.Total
	paid := int64(0)

	for _, playerID := range orderedPlayers {
		if remaining <= 0 {
			return payouts
		}

		payouts[playerID] = 0
		playerAmount := p.PlayerTotals[playerID] - paid
		paid += playerAmount
		for _, otherPlayerTotal := range p.PlayerTotals {
			amount := utils.MinInt64(otherPlayerTotal, playerAmount)
			remaining -= amount
			payouts[playerID] += amount
		}

		delete(p.PlayerTotals, playerID)
	}

	return payouts
}
