package holdem

import (
	"qpoker/cards"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameDeal(t *testing.T) {
	cases := [][]*Player{
		[]*Player{
			&Player{ID: 7, Stack: int64(100)},
			&Player{ID: 8, Stack: int64(100)},
		},
		[]*Player{
			&Player{ID: 2, Stack: int64(100)},
			&Player{ID: 3, Stack: int64(100)},
			&Player{ID: 4, Stack: int64(100)},
			&Player{ID: 5, Stack: int64(100)},
			&Player{ID: 6, Stack: int64(100)},
			&Player{ID: 7, Stack: int64(100)},
			&Player{ID: 8, Stack: int64(100)},
		},
	}

	for _, players := range cases {
		holdEm := NewHoldEm(NewTable(len(players), players))
		assert.Equal(t, holdEm.State, StateInit)

		holdEm.Deal()
		assert.Equal(t, holdEm.State, StateDeal)

		activePlayers := holdEm.Table.GetActivePlayers()
		assert.Equal(t, len(players), len(activePlayers))

		for _, player := range activePlayers {
			assert.Equal(t, 2, len(player.Cards))
		}
	}
}

func TestGameAdvance(t *testing.T) {
	players := NewTable(9, []*Player{
		&Player{ID: 5, Stack: int64(100)},
		&Player{ID: 6, Stack: int64(100)},
		&Player{ID: 7, Stack: int64(100)},
		&Player{ID: 8, Stack: int64(100)},
	})

	holdEm := NewHoldEm(players)
	holdEm.Deal()

	// Flop
	err := holdEm.Advance()
	assert.NoError(t, err)
	assert.Equal(t, StateFlop, holdEm.State)
	assert.Equal(t, 3, len(holdEm.Board))

	// Turn
	err = holdEm.Advance()
	assert.NoError(t, err)
	assert.Equal(t, StateTurn, holdEm.State)
	assert.Equal(t, 4, len(holdEm.Board))

	// River
	err = holdEm.Advance()
	assert.NoError(t, err)
	assert.Equal(t, StateRiver, holdEm.State)
	assert.Equal(t, 5, len(holdEm.Board))

	// Error, nowhere left to go
	err = holdEm.Advance()
	assert.Error(t, err)
}

func TestGameSimulate(t *testing.T) {
	table := NewTable(12, []*Player{
		&Player{ID: 5, Stack: int64(2)},
		&Player{ID: 6, Stack: int64(2)},
		&Player{ID: 7, Stack: int64(2)},
		&Player{ID: 8, Stack: int64(2)},
	})

	holdEm := NewHoldEm(table)

	holdEm.Deal()
	err := holdEm.Simulate()

	assert.NoError(t, err)
	assert.Equal(t, StateRiver, holdEm.State)
	assert.Equal(t, 5, len(holdEm.Board))
}

func TestGetWinningIDs(t *testing.T) {
	type TestCase struct {
		state    *HoldEm
		expected [][]int64
	}
	testCases := []TestCase{
		TestCase{
			state: &HoldEm{
				Table: NewTable(12, []*Player{
					&Player{
						ID: 5,
						Cards: []cards.Card{
							cards.NewCard(1, cards.SuitClubs),
							cards.NewCard(1, cards.SuitHearts),
						},
					},
					&Player{
						ID: 7,
						Cards: []cards.Card{
							cards.NewCard(13, cards.SuitClubs),
							cards.NewCard(13, cards.SuitHearts),
						},
					},
				}),
				Board: []cards.Card{
					cards.NewCard(1, cards.SuitDiamonds),
					cards.NewCard(13, cards.SuitDiamonds),
					cards.NewCard(5, cards.SuitDiamonds),
					cards.NewCard(2, cards.SuitSpades),
					cards.NewCard(3, cards.SuitSpades),
				},
			},
			expected: [][]int64{[]int64{5}, []int64{7}},
		},
		TestCase{
			state: &HoldEm{
				Table: NewTable(12, []*Player{
					&Player{
						ID: 1,
						Cards: []cards.Card{
							cards.NewCard(1, cards.SuitClubs),
							cards.NewCard(1, cards.SuitHearts),
						},
					},
					&Player{
						ID: 2,
						Cards: []cards.Card{
							cards.NewCard(9, cards.SuitClubs),
							cards.NewCard(8, cards.SuitHearts),
						},
					},
					&Player{
						ID: 3,
						Cards: []cards.Card{
							cards.NewCard(11, cards.SuitClubs),
							cards.NewCard(11, cards.SuitHearts),
						},
					},
					&Player{
						ID: 4,
						Cards: []cards.Card{
							cards.NewCard(11, cards.SuitDiamonds),
							cards.NewCard(11, cards.SuitSpades),
						},
					},
					&Player{
						ID: 5,
						Cards: []cards.Card{
							cards.NewCard(8, cards.SuitClubs),
							cards.NewCard(9, cards.SuitHearts),
						},
					},
				}),
				Board: []cards.Card{
					cards.NewCard(1, cards.SuitSpades),
					cards.NewCard(13, cards.SuitSpades),
					cards.NewCard(5, cards.SuitDiamonds),
					cards.NewCard(2, cards.SuitDiamonds),
					cards.NewCard(3, cards.SuitDiamonds),
				},
			},
			expected: [][]int64{[]int64{1}, []int64{4, 3}, []int64{5, 2}},
		},
	}

	for _, c := range testCases {
		assert.Equal(t, c.expected, c.state.GetWinningIDs())
	}
}
