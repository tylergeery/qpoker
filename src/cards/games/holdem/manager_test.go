package holdem

import (
	"fmt"
	"qpoker/models"
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
		_, err := NewGameManager(c.players, models.GameOptions{})
		assert.Equal(t, c.expected, err.Error())
	}
}

func TestPlayTooManyPlayers(t *testing.T) {
	players := []*Player{}
	for i := 0; i < 8; i++ {
		players = append(players, &Player{ID: int64(i)})
	}

	_, err := NewGameManager(players, models.GameOptions{Capacity: 4, BigBlind: 100})
	assert.Error(t, err)
}

func TestPlayHandNotEnoughPlayersDueToStacks(t *testing.T) {
	players := []*Player{}
	for i := 0; i < 5; i++ {
		players = append(players, &Player{ID: int64(i)})
	}

	gm, err := NewGameManager(players, models.GameOptions{Capacity: 5, BigBlind: 100})
	assert.NoError(t, err)

	err = gm.NextHand()
	assert.Error(t, err)
}

func TestPlayHandAllFold(t *testing.T) {
	var nilMap map[string]bool

	players := []*Player{}
	for i := 0; i < 5; i++ {
		players = append(players, &Player{ID: int64(i), Stack: int64(200)})
	}

	gm, err := NewGameManager(players, models.GameOptions{Capacity: 5, BigBlind: 100})
	assert.NoError(t, err)

	// 1 becomes dealer, 2 LB, 3 BB, 4 active
	err = gm.NextHand()
	assert.NoError(t, err)
	assert.Equal(t, int64(100), players[3].Stack)
	assert.Equal(t, int64(100), gm.Pot.PlayerBets[players[3].ID], fmt.Sprintf("%+v", gm.Pot.PlayerBets))
	assert.Equal(t, int64(150), players[2].Stack)
	assert.Equal(t, int64(50), gm.Pot.PlayerBets[players[2].ID], fmt.Sprintf("%+v", gm.Pot.PlayerBets))

	// Try to move out of turn
	for i, player := range players {
		if i != 4 {
			assert.Equal(t, nilMap, player.Options)
			_, err = gm.PlayerAction(player.ID, NewActionCall())
			assert.Error(t, err)
		}
	}

	assert.True(t, players[4].Options["can_bet"])
	assert.True(t, players[4].Options["can_call"])
	assert.True(t, players[4].Options["can_fold"])
	assert.False(t, players[4].Options["can_check"])
	complete, err := gm.PlayerAction(players[4].ID, NewActionFold())
	assert.NoError(t, err)
	assert.False(t, complete)
	assert.Equal(t, nilMap, players[4].Options)

	assert.True(t, players[0].Options["can_bet"])
	assert.True(t, players[0].Options["can_call"])
	assert.True(t, players[0].Options["can_fold"])
	assert.False(t, players[0].Options["can_check"])
	complete, err = gm.PlayerAction(players[0].ID, NewActionFold())
	assert.NoError(t, err)
	assert.False(t, complete)
	assert.Equal(t, nilMap, players[0].Options)

	complete, err = gm.PlayerAction(players[1].ID, NewActionFold())
	assert.NoError(t, err)
	assert.False(t, complete)
	complete, err = gm.PlayerAction(players[2].ID, NewActionFold())
	assert.NoError(t, err)
	assert.True(t, complete)
}

func TestPlayHandAllCheckAndCall(t *testing.T) {
	players := []*Player{}
	for i := 0; i < 3; i++ {
		players = append(players, &Player{ID: int64(i), Stack: int64(200)})
	}

	gm, err := NewGameManager(players, models.GameOptions{Capacity: 5, BigBlind: 100})
	assert.NoError(t, err)

	// 1 becomes dealer, 2 LB, 0 BB, 1 active
	err = gm.NextHand()
	assert.NoError(t, err)

	complete, err := gm.PlayerAction(players[1].ID, NewActionCall())
	assert.NoError(t, err)
	assert.False(t, complete)

	complete, err = gm.PlayerAction(players[2].ID, NewActionCall())
	assert.NoError(t, err)
	assert.False(t, complete)

	assert.True(t, players[0].Options["can_bet"])
	assert.False(t, players[0].Options["can_call"])
	assert.False(t, players[0].Options["can_fold"])
	assert.True(t, players[0].Options["can_check"])
	complete, err = gm.PlayerAction(players[0].ID, NewActionCheck())
	assert.NoError(t, err)
	assert.False(t, complete)
	assert.Equal(t, gm.State.State, StateFlop)

	assert.True(t, players[2].Options["can_bet"])
	assert.False(t, players[2].Options["can_call"])
	assert.False(t, players[2].Options["can_fold"])
	assert.True(t, players[2].Options["can_check"])
	complete, err = gm.PlayerAction(players[2].ID, NewActionCheck())
	assert.NoError(t, err)
	assert.False(t, complete)

	assert.True(t, players[0].Options["can_bet"])
	assert.False(t, players[0].Options["can_call"])
	assert.False(t, players[0].Options["can_fold"])
	assert.True(t, players[0].Options["can_check"])
	complete, err = gm.PlayerAction(players[0].ID, NewActionCheck())
	assert.NoError(t, err)
	assert.False(t, complete)

	assert.True(t, players[1].Options["can_bet"])
	assert.False(t, players[1].Options["can_call"])
	assert.False(t, players[1].Options["can_fold"])
	assert.True(t, players[1].Options["can_check"])
	complete, err = gm.PlayerAction(players[1].ID, NewActionCheck())
	assert.NoError(t, err)
	assert.False(t, complete)
	assert.Equal(t, gm.State.State, StateTurn)

	complete, err = gm.PlayerAction(players[2].ID, NewActionCheck())
	assert.NoError(t, err)
	assert.False(t, complete)

	complete, err = gm.PlayerAction(players[0].ID, NewActionCheck())
	assert.NoError(t, err)
	assert.False(t, complete)

	complete, err = gm.PlayerAction(players[1].ID, NewActionCheck())
	assert.NoError(t, err)
	assert.False(t, complete)
	assert.Equal(t, gm.State.State, StateRiver)

	complete, err = gm.PlayerAction(players[2].ID, NewActionCheck())
	assert.NoError(t, err)
	assert.False(t, complete)

	complete, err = gm.PlayerAction(players[0].ID, NewActionCheck())
	assert.NoError(t, err)
	assert.False(t, complete)

	complete, err = gm.PlayerAction(players[1].ID, NewActionCheck())
	assert.NoError(t, err)
	assert.True(t, complete)
	assert.Equal(t, gm.State.State, StateRiver)
}

func TestPlayComplexHand(t *testing.T) {

}

func TestPlayManyHands(t *testing.T) {

}
