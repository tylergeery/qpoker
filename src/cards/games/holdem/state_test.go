package holdem

import (
	"qpoker/cards/games"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGameError(t *testing.T) {
	type TestCase struct {
		players  []*games.Player
		expected string
	}
	cases := []TestCase{
		TestCase{
			players:  []*games.Player{},
			expected: "Invalid player count: 0",
		},
		TestCase{
			players: []*games.Player{
				&games.Player{ID: 1},
			},
			expected: "Invalid player count: 1",
		},
		TestCase{
			players: []*games.Player{
				&games.Player{ID: 1},
				&games.Player{ID: 2},
				&games.Player{ID: 3},
				&games.Player{ID: 4},
				&games.Player{ID: 5},
				&games.Player{ID: 6},
				&games.Player{ID: 7},
				&games.Player{ID: 8},
				&games.Player{ID: 9},
				&games.Player{ID: 10},
				&games.Player{ID: 11},
				&games.Player{ID: 12},
				&games.Player{ID: 13},
			},
			expected: "Invalid player count: 13",
		},
	}

	for _, c := range cases {
		_, err := NewHoldEm(c.players)
		assert.Equal(t, c.expected, err.Error())
	}
}

func TestGameDeal(t *testing.T) {
	cases := [][]*games.Player{
		[]*games.Player{
			&games.Player{ID: 7},
			&games.Player{ID: 8},
		},
		[]*games.Player{
			&games.Player{ID: 2},
			&games.Player{ID: 3},
			&games.Player{ID: 4},
			&games.Player{ID: 5},
			&games.Player{ID: 6},
			&games.Player{ID: 7},
			&games.Player{ID: 8},
		},
	}

	for _, players := range cases {
		holdEm, err := NewHoldEm(players)
		assert.NoError(t, err)
		assert.Equal(t, holdEm.State, StateInit)

		holdEm.Deal()
		assert.Equal(t, holdEm.State, StateDeal)

		for _, player := range holdEm.Players {
			assert.Equal(t, 2, len(player.Cards))
		}
	}
}

func TestGameAdvance(t *testing.T) {
	players := []*games.Player{
		&games.Player{ID: 5},
		&games.Player{ID: 6},
		&games.Player{ID: 7},
		&games.Player{ID: 8},
	}

	holdEm, err := NewHoldEm(players)
	assert.NoError(t, err)

	holdEm.Deal()

	// Flop
	err = holdEm.Advance()
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
	players := []*games.Player{
		&games.Player{ID: 5},
		&games.Player{ID: 6},
		&games.Player{ID: 7},
		&games.Player{ID: 8},
	}

	holdEm, err := NewHoldEm(players)
	assert.NoError(t, err)

	holdEm.Deal()
	err = holdEm.Simulate()

	assert.NoError(t, err)
	assert.Equal(t, StateRiver, holdEm.State)
	assert.Equal(t, 5, len(holdEm.Board))
}
