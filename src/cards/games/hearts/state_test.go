package holdem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameDeal(t *testing.T) {
	cases := [][]*Player{
		[]*Player{
			&Player{ID: 7, Score: 0},
			&Player{ID: 8, Score: 0},
		},
		[]*Player{
			&Player{ID: 2, Score: 0},
			&Player{ID: 3, Score: 0},
			&Player{ID: 4, Score: 0},
			&Player{ID: 5, Score: 0},
			&Player{ID: 6, Score: 0},
			&Player{ID: 7, Score: 0},
			&Player{ID: 8, Score: 0},
		},
	}

	for _, players := range cases {
		hearts := NewHearts(NewTable(4, players), StatePassing)
		assert.Equal(t, hearts.State, StatePassing)

		hearts.Deal()
		assert.Equal(t, hearts.State, StateActive)

		activePlayers := hearts.Table.GetActivePlayers()
		assert.Equal(t, len(players), len(activePlayers))

		for _, player := range activePlayers {
			assert.Equal(t, 13, len(player.Cards))
		}
	}
}
