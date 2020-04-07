package holdem

import (
	"qpoker/cards/games"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplePot(t *testing.T) {
	// Given
	player1 := &games.Player{ID: 1}
	player2 := &games.Player{ID: 2}
	player3 := &games.Player{ID: 3}

	pot := NewPot([]*games.Player{player1, player2, player3})
	pot.AddBet(player1.ID, 500)
	pot.AddBet(player2.ID, 500)
	pot.AddBet(player3.ID, 500)

	// When
	payouts := pot.GetPayouts([]int64{player3.ID, player2.ID, player1.ID})

	// Then
	_, ok1 := payouts[player1.ID]
	_, ok2 := payouts[player2.ID]
	payout3, _ := payouts[player3.ID]

	assert.False(t, ok1)
	assert.False(t, ok2)
	assert.Equal(t, int64(1500), payout3)
}

func TestComplexPot(t *testing.T) {
	// Given
	player1 := &games.Player{ID: 1}
	player2 := &games.Player{ID: 2}
	player3 := &games.Player{ID: 3}
	player4 := &games.Player{ID: 4}

	pot := NewPot([]*games.Player{player1, player2, player3, player4})
	pot.AddBet(player4.ID, 100)
	pot.AddBet(player1.ID, 500)
	pot.AddBet(player2.ID, 500)
	pot.AddBet(player3.ID, 1000)
	pot.AddBet(player4.ID, 822) // all in
	pot.AddBet(player1.ID, 133) // all in
	pot.AddBet(player2.ID, 500)
	pot.ClearBets()

	pot.AddBet(player3.ID, 1000)
	pot.AddBet(player2.ID, 1000)

	// When
	payouts := pot.GetPayouts([]int64{player1.ID, player4.ID, player3.ID, player2.ID})

	// Then
	payout1 := payouts[player1.ID]
	_, ok2 := payouts[player2.ID]
	payout3 := payouts[player3.ID]
	payout4 := payouts[player4.ID]

	assert.False(t, ok2)
	assert.Equal(t, int64(2532), payout1)
	assert.Equal(t, int64(2156), payout3)
	assert.Equal(t, int64(867), payout4)

}
