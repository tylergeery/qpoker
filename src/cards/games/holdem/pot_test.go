package holdem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplePot(t *testing.T) {
	// Given
	player1 := &Player{ID: 1}
	player2 := &Player{ID: 2}
	player3 := &Player{ID: 3}

	pot := NewPot([]*Player{player1, player2, player3})
	pot.AddBet(player1.ID, 500)
	pot.AddBet(player2.ID, 500)
	pot.AddBet(player3.ID, 500)

	// When
	payouts := pot.GetPayouts([][]int64{
		[]int64{player3.ID},
		[]int64{player2.ID},
		[]int64{player1.ID},
	})

	// Then
	payout1 := payouts[player1.ID]
	payout2 := payouts[player2.ID]
	payout3 := payouts[player3.ID]

	assert.Equal(t, int64(0), payout1)
	assert.Equal(t, int64(0), payout2)
	assert.Equal(t, int64(1500), payout3)
}

func TestComplexPot(t *testing.T) {
	// Given
	player1 := &Player{ID: 1}
	player2 := &Player{ID: 2}
	player3 := &Player{ID: 3}
	player4 := &Player{ID: 4}

	pot := NewPot([]*Player{player1, player2, player3, player4})
	pot.AddBet(player4.ID, 100)
	pot.AddBet(player1.ID, 500)
	pot.AddBet(player2.ID, 500)

	assert.Equal(t, pot.MaxBet(), int64(500))

	pot.AddBet(player3.ID, 1000)

	assert.Equal(t, pot.MaxBet(), int64(1000))

	pot.AddBet(player4.ID, 822) // all in
	pot.AddBet(player1.ID, 133) // all in
	pot.AddBet(player2.ID, 500)

	assert.Equal(t, pot.MaxBet(), int64(1000))

	pot.ClearBets()
	assert.Equal(t, pot.MaxBet(), int64(0))

	pot.AddBet(player3.ID, 1000)
	pot.AddBet(player2.ID, 1000)

	// When
	payouts := pot.GetPayouts([][]int64{
		[]int64{player1.ID}, // bet 633
		[]int64{player4.ID}, // bet 922
		[]int64{player3.ID}, // bet 2000
		[]int64{player2.ID}, // bet 2000
	})

	// Then
	payout1 := payouts[player1.ID]
	payout2 := payouts[player2.ID]
	payout3 := payouts[player3.ID]
	payout4 := payouts[player4.ID]

	assert.Equal(t, int64(0), payout2)
	assert.Equal(t, int64(2532), payout1)
	assert.Equal(t, int64(2156), payout3)
	assert.Equal(t, int64(867), payout4)

}

func TestGetPayouts(t *testing.T) {
	type TestCase struct {
		pot             *Pot
		orderedPlayers  [][]int64
		expectedPayouts map[int64]int64
	}
	testCases := []TestCase{
		TestCase{
			pot: &Pot{
				PlayerBets: map[int64]int64{},
				PlayerTotals: map[int64]int64{
					1: 50,
					2: 50,
				},
				Total: 100,
			},
			orderedPlayers: [][]int64{
				[]int64{1}, []int64{2},
			},
			expectedPayouts: map[int64]int64{
				1: 100,
				2: 0,
			},
		},
		TestCase{
			pot: &Pot{
				PlayerBets: map[int64]int64{},
				PlayerTotals: map[int64]int64{
					1: 4921,
					2: 4921,
				},
				Total: 9842,
			},
			orderedPlayers: [][]int64{
				[]int64{1, 2}, []int64{},
			},
			expectedPayouts: map[int64]int64{
				1: 4921,
				2: 4921,
			},
		},
		TestCase{
			pot: &Pot{
				PlayerBets: map[int64]int64{},
				PlayerTotals: map[int64]int64{
					1: 50,
					2: 50,
					3: 30,
				},
				Total: 130,
			},
			orderedPlayers: [][]int64{
				[]int64{1, 3}, []int64{2},
			},
			expectedPayouts: map[int64]int64{
				1: 85,
				2: 0,
				3: 45,
			},
		},
		TestCase{
			pot: &Pot{
				PlayerBets: map[int64]int64{},
				PlayerTotals: map[int64]int64{
					1: 5020,
					2: 540,
					3: 780,
					4: 200,
				},
				Total: 6540,
			},
			orderedPlayers: [][]int64{
				[]int64{3}, []int64{4, 1}, []int64{2},
			},
			expectedPayouts: map[int64]int64{
				3: 2300,
				1: 4240,
				4: 0,
				2: 0,
			},
		},
		TestCase{
			pot: &Pot{
				PlayerBets: map[int64]int64{},
				PlayerTotals: map[int64]int64{
					1: 5400,
					2: 5400,
					3: 780,
					4: 2000,
					5: 8000,
					6: 8000,
				},
				Total: 29580,
			},
			orderedPlayers: [][]int64{
				[]int64{4}, []int64{2, 3}, []int64{1},
				[]int64{5}, []int64{6},
			},
			expectedPayouts: map[int64]int64{
				4: 10780,
				3: 0,
				2: 13600,
				1: 0,
				5: 5200,
				6: 0,
			},
		},
	}

	for _, c := range testCases {
		payouts := c.pot.GetPayouts(c.orderedPlayers)
		assert.Equal(t, c.expectedPayouts, payouts)
	}
}
