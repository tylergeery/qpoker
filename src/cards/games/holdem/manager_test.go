package holdem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGameError(t *testing.T) {
	type TestCase struct {
		players  []*Player
		expected string
	}
	cases := []TestCase{
		TestCase{
			players:  []*Player{},
			expected: "Invalid player count: 0",
		},
		TestCase{
			players: []*Player{
				&Player{ID: 1},
			},
			expected: "Invalid player count: 1",
		},
		TestCase{
			players: []*Player{
				&Player{ID: 1},
				&Player{ID: 2},
				&Player{ID: 3},
				&Player{ID: 4},
				&Player{ID: 5},
				&Player{ID: 6},
				&Player{ID: 7},
				&Player{ID: 8},
				&Player{ID: 9},
				&Player{ID: 10},
				&Player{ID: 11},
				&Player{ID: 12},
				&Player{ID: 13},
			},
			expected: "Invalid player count: 13",
		},
	}

	for _, c := range cases {
		_, err := NewGameManager(c.players, GameOptions{})
		assert.Equal(t, c.expected, err.Error())
	}
}

func TestPlayTooManyPlayers(t *testing.T) {
	players := []*Player{}
	for i := 0; i < 8; i++ {
		players = append(players, &Player{ID: int64(i)})
	}

	_, err := NewGameManager(players, GameOptions{Capacity: 4, BigBlind: 100})
	assert.Error(t, err)
}

func TestPlayHandNotEnoughPlayersDueToStacks(t *testing.T) {
	players := []*Player{}
	for i := 0; i < 5; i++ {
		players = append(players, &Player{ID: int64(i)})
	}

	gm, err := NewGameManager(players, GameOptions{Capacity: 5, BigBlind: 100})
	assert.NoError(t, err)

	err = gm.NextHand()
	assert.Error(t, err)
}

func TestPlayHandAllFold(t *testing.T) {
	players := []*Player{}
	for i := 0; i < 5; i++ {
		players = append(players, &Player{ID: 0, Stack: int64(200)})
	}

	_, err := NewGameManager(players, GameOptions{Capacity: 5, BigBlind: 100})
	assert.NoError(t, err)
}

func TestPlayHandAllBetAndCall(t *testing.T) {

}

func TestPlayComplexHand(t *testing.T) {

}

func TestPlayManyHands(t *testing.T) {

}
