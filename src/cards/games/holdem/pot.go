package holdem

import (
	"fmt"
	"qpoker/qutils"
)

// Pot controls betting pot
type Pot struct {
	Payouts      map[int64]int64 `json:"payouts"`
	PlayerBets   map[int64]int64 `json:"player_bets"`
	PlayerTotals map[int64]int64 `json:"player_totals"`
	Total        int64           `json:"total"`
}

// NewPot returns a new Pot instance for the current group of players
func NewPot(players []*Player) *Pot {
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

// MaxBet returns the max current bet
func (p *Pot) MaxBet() int64 {
	max := int64(0)

	for _, bet := range p.PlayerBets {
		if bet > max {
			max = bet
		}
	}

	return max
}

func (p *Pot) minTotalAmongstPlayers(playerIDs []int64) int64 {
	amounts := []int64{}

	for _, playerID := range playerIDs {
		amounts = append(amounts, p.PlayerTotals[playerID])
	}

	return qutils.MinInt64(amounts...)
}

func (p *Pot) playersWithMoreThanTotal(amount int64, playerIDs []int64) []int64 {
	remainingPlayerIDs := []int64{}

	for _, playerID := range playerIDs {
		if p.PlayerTotals[playerID] > amount {
			remainingPlayerIDs = append(remainingPlayerIDs, playerID)
		}
	}

	return remainingPlayerIDs
}

func (p *Pot) splitEvenlyAmongstPlayers(amount int64, playerIDs []int64, payouts map[int64]int64) {
	portion := float64(1) / float64(len(playerIDs))
	perPlayer := int64(portion * float64(amount))
	remainder := amount - (perPlayer * int64(len(playerIDs)))

	for _, playerID := range playerIDs {
		payouts[playerID] += perPlayer

		// handle rounding errors so game doesnt swallow money
		if remainder > 0 {
			remainder--
			payouts[playerID]++
		}
	}
}

func (p *Pot) divyUpPlayerTotal(playerID, amount, paid int64, orderedPlayerIDs [][]int64, payouts map[int64]int64) {
	if amount <= 0 {
		return
	}

	if len(orderedPlayerIDs) == 0 {
		fmt.Printf("Error: %d unpaid from player %d, payouts: %+v\n", amount, playerID, payouts)
		return
	}

	// if user won the hand outright, give them back their bet
	if paid == 0 && qutils.Int64SliceHasValue(orderedPlayerIDs[0], playerID) {
		payouts[playerID] += amount
		return
	}

	// only consider players who bet more than what has already been paid
	nextTier, prevTier := paid, paid
	remainingPlayerIDs := orderedPlayerIDs[0]
	for nextTier < amount {
		remainingPlayerIDs := p.playersWithMoreThanTotal(nextTier, remainingPlayerIDs)
		if len(remainingPlayerIDs) == 0 {
			// remaining amount will go to next best hands
			p.divyUpPlayerTotal(playerID, amount, nextTier, orderedPlayerIDs[1:], payouts)
			return
		}

		prevTier, nextTier = nextTier, qutils.MinInt64(amount, p.minTotalAmongstPlayers(remainingPlayerIDs))
		p.splitEvenlyAmongstPlayers(nextTier-prevTier, remainingPlayerIDs, payouts)
	}
}

func (p *Pot) initPayouts() map[int64]int64 {
	payouts := map[int64]int64{}
	for playerID := range p.PlayerTotals {
		payouts[playerID] = 0
	}

	return payouts
}

// GetPayouts returns the winning amounts for each user by winning priority
func (p *Pot) GetPayouts(orderedPlayers [][]int64) map[int64]int64 {
	p.ClearBets() // put everything to the pot

	payouts := p.initPayouts()

	// loop through all players
	for playerID, playerTotal := range p.PlayerTotals {
		p.divyUpPlayerTotal(playerID, playerTotal, 0, orderedPlayers, payouts)
	}

	p.Payouts = payouts

	return payouts
}
