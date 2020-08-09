package hearts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameDeal(t *testing.T) {
	cases := [][]*Player{
		[]*Player{
			&Player{ID: 7, Score: 0},
			&Player{ID: 8, Score: 0},
			&Player{ID: 2, Score: 0},
			&Player{ID: 3, Score: 0},
		},
		[]*Player{
			&Player{ID: 4, Score: 0},
			&Player{ID: 5, Score: 0},
			&Player{ID: 6, Score: 0},
			&Player{ID: 7, Score: 0},
		},
	}

	for _, players := range cases {
		hearts := NewHearts(NewTable(4, players), StatePassing)
		hearts.Deal()
		assert.Equal(t, hearts.State, StatePassing)

		activePlayers := hearts.Table.GetActivePlayers()
		assert.Equal(t, len(players), len(activePlayers))

		for _, player := range activePlayers {
			hearts.addPass(player.ID, player.Cards[10:])
		}

		assert.Equal(t, hearts.State, StateActive)
		assert.Equal(t, len(players), len(activePlayers))

		for _, player := range activePlayers {
			assert.Equal(t, 13, len(player.Cards), fmt.Sprintf("Cards: %+v", player.Cards))
		}
	}
}
