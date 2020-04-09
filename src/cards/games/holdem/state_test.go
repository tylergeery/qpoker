package holdem

import (
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
	players := NewTable(12, []*Player{
		&Player{ID: 5, Stack: int64(2)},
		&Player{ID: 6, Stack: int64(2)},
		&Player{ID: 7, Stack: int64(2)},
		&Player{ID: 8, Stack: int64(2)},
	})

	holdEm := NewHoldEm(players)

	holdEm.Deal()
	err := holdEm.Simulate()

	assert.NoError(t, err)
	assert.Equal(t, StateRiver, holdEm.State)
	assert.Equal(t, 5, len(holdEm.Board))
}
