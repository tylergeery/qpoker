package holdem

import (
	"fmt"
	"qpoker/models"
	"qpoker/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestPlayer(t *testing.T, stack int) *Player {
	player := models.CreateTestPlayer()

	return &Player{
		ID:       player.ID,
		Username: player.Username,
		State:    PlayerStateInit,
		Stack:    int64(stack),
	}
}

func createTestManager(t *testing.T, players ...*Player) *GameManager {
	game := models.CreateTestGame(players[0].ID)

	manager, err := NewGameManager(game.ID, players, models.GameOptions{BigBlind: 50})
	assert.NoError(t, err)
	err = manager.NextHand()
	assert.NoError(t, err)

	return manager
}

func expectedPlayerActions(canBet, canCall, canCheck, canFold bool) map[string]bool {
	return map[string]bool{
		"can_bet":   canBet,
		"can_call":  canCall,
		"can_check": canCheck,
		"can_fold":  canFold,
	}
}

func TestNewGameError(t *testing.T) {
	type TestCase struct {
		players  []*Player
		expected string
	}
	cases := []TestCase{
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
		_, err := NewGameManager(int64(0), c.players, models.GameOptions{})
		assert.Equal(t, c.expected, err.Error())
	}
}

func TestPlayTooManyPlayers(t *testing.T) {
	players := []*Player{}
	for i := 0; i < 8; i++ {
		players = append(players, &Player{ID: int64(i)})
	}

	_, err := NewGameManager(int64(0), players, models.GameOptions{Capacity: 4, BigBlind: 100})
	assert.Error(t, err)
}

func TestPlayHandNotEnoughPlayersDueToStacks(t *testing.T) {
	players := []*Player{}
	for i := 0; i < 5; i++ {
		players = append(players, &Player{ID: int64(i)})
	}

	gm, err := NewGameManager(int64(0), players, models.GameOptions{Capacity: 5, BigBlind: 100})
	assert.NoError(t, err)

	err = gm.NextHand()
	assert.Error(t, err)
}

func TestPlayHandAllFold(t *testing.T) {
	var nilMap map[string]bool
	var player *models.Player

	players := []*Player{}
	for i := 0; i < 5; i++ {
		player = models.CreateTestPlayer()
		players = append(players, &Player{ID: player.ID, Stack: int64(200)})
	}
	game := models.CreateTestGame(player.ID)

	gm, err := NewGameManager(game.ID, players, models.GameOptions{Capacity: 5, BigBlind: 100})
	assert.NoError(t, err)

	// 1 becomes dealer, 2 LB, 3 BB, 4 active
	err = gm.NextHand()
	assert.NoError(t, err)
	assert.Equal(t, int64(100), players[3].Stack)
	assert.Equal(t, int64(100), gm.Pot.PlayerBets[players[3].ID], fmt.Sprintf("%+v", gm.Pot))
	assert.Equal(t, int64(150), players[2].Stack)
	assert.Equal(t, int64(50), gm.Pot.PlayerBets[players[2].ID], fmt.Sprintf("%+v", gm.Pot.PlayerBets))

	// check game hand saved expected values
	assert.Greater(t, gm.gameHand.ID, int64(0))
	assert.Equal(t, gm.gameHand.GameID, game.ID)
	for i := 0; i < 5; i++ {
		assert.Greater(t, gm.gamePlayerHands[players[i].ID].ID, int64(0))
		assert.Equal(t, gm.gamePlayerHands[players[i].ID].GameHandID, gm.gameHand.ID)
		assert.Equal(t, players[i].ID, gm.gamePlayerHands[players[i].ID].PlayerID)
		assert.Equal(t, 2, len(gm.gamePlayerHands[players[i].ID].Cards))
		assert.Equal(t, int64(200), gm.gamePlayerHands[players[i].ID].StartingStack)
	}

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

	// check game hand final save expected values
	expectedFinalStacks := []int64{200, 200, 150, 250, 200}
	assert.Equal(t, 0, len(gm.gameHand.Board))
	assert.Equal(t, int64(150), gm.gameHand.Payouts[players[3].ID])
	assert.Equal(t, int64(100), gm.gameHand.Bets[players[3].ID])
	assert.Equal(t, int64(50), gm.gameHand.Bets[players[2].ID])
	assert.Equal(t, gm.gameHand.GameID, game.ID)
	for i := 0; i < 5; i++ {
		assert.Greater(t, gm.gamePlayerHands[players[i].ID].ID, int64(0))
		assert.Equal(t, 2, len(gm.gamePlayerHands[players[i].ID].Cards))
		assert.Equal(t, int64(200), gm.gamePlayerHands[players[i].ID].StartingStack)
		assert.Equal(t, expectedFinalStacks[i], gm.gamePlayerHands[players[i].ID].EndingStack, fmt.Sprintf("Final Player Hand: pos(%d) %+v\n", i, gm.gamePlayerHands[players[i].ID]))
		assert.False(t, players[i].CardsVisible)
	}
}

func TestPlayHandAllCheckAndCall(t *testing.T) {
	var player *models.Player
	players := []*Player{}
	for i := 0; i < 3; i++ {
		player = models.CreateTestPlayer()
		players = append(players, &Player{ID: player.ID, Stack: int64(200)})
	}
	game := models.CreateTestGame(player.ID)

	gm, err := NewGameManager(game.ID, players, models.GameOptions{Capacity: 5, BigBlind: 100})
	assert.NoError(t, err)

	// 1 becomes dealer, 2 LB, 0 BB, 1 active
	err = gm.NextHand()
	assert.NoError(t, err)

	// check game hand saved expected values
	assert.Greater(t, gm.gameHand.ID, int64(0))
	assert.Equal(t, gm.gameHand.GameID, game.ID)
	for i := 0; i < 3; i++ {
		assert.Greater(t, gm.gamePlayerHands[players[i].ID].ID, int64(0))
		assert.Equal(t, gm.gamePlayerHands[players[i].ID].GameHandID, gm.gameHand.ID)
		assert.Equal(t, players[i].ID, gm.gamePlayerHands[players[i].ID].PlayerID)
		assert.Equal(t, 2, len(gm.gamePlayerHands[players[i].ID].Cards))
		assert.Equal(t, int64(200), gm.gamePlayerHands[players[i].ID].StartingStack)
	}

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

	// check game hand final save expected values
	min, max, expectedMin, expectedMax := int64(500), int64(0), int64(100), int64(400)
	total, expectedTotal := int64(0), int64(3*200)
	assert.Equal(t, 5, len(gm.gameHand.Board))
	assert.Equal(t, gm.gameHand.GameID, game.ID)
	for i := 0; i < 3; i++ {
		assert.Greater(t, gm.gamePlayerHands[players[i].ID].ID, int64(0))
		assert.Equal(t, 2, len(gm.gamePlayerHands[players[i].ID].Cards))
		assert.Equal(t, int64(200), gm.gamePlayerHands[players[i].ID].StartingStack)
		assert.Equal(t, int64(100), gm.gameHand.Bets[players[i].ID])
		min, max = utils.MinInt64(min, gm.gamePlayerHands[players[i].ID].EndingStack), utils.MaxInt64(max, gm.gamePlayerHands[players[i].ID].EndingStack)
		total += gm.gamePlayerHands[players[i].ID].EndingStack
		assert.True(t, players[i].CardsVisible)
	}

	assert.Equal(t, expectedMin, min)
	assert.Equal(t, expectedMax, max)
	assert.Equal(t, expectedTotal, total)
}

func TestGetPlayerActions(t *testing.T) {
	type TestCase struct {
		manager  *GameManager
		expected map[string]bool
	}
	gamePlayers := []*Player{
		createTestPlayer(t, 1000),
		createTestPlayer(t, 1000),
		createTestPlayer(t, 1000),
	}
	advancedManager := createTestManager(t, gamePlayers...)
	complete, err := advancedManager.PlayerAction(gamePlayers[1].ID, Action{ActionCall, int64(0)})
	assert.False(t, complete)
	assert.NoError(t, err)
	complete, err = advancedManager.PlayerAction(gamePlayers[2].ID, Action{ActionCall, int64(0)})
	assert.False(t, complete)
	assert.NoError(t, err)

	cases := []TestCase{
		TestCase{
			manager: createTestManager(t,
				createTestPlayer(t, 25),
				createTestPlayer(t, 50)),
			expected: map[string]bool{
				"can_bet":   false,
				"can_call":  false,
				"can_check": false,
				"can_fold":  false,
			},
		},
		TestCase{
			manager: createTestManager(t,
				createTestPlayer(t, 50),
				createTestPlayer(t, 50),
				createTestPlayer(t, 50)),
			expected: map[string]bool{
				"can_bet":   false,
				"can_call":  true,
				"can_check": false,
				"can_fold":  true,
			},
		},
		TestCase{
			manager: advancedManager,
			expected: map[string]bool{
				"can_bet":   true,
				"can_call":  false,
				"can_check": true,
				"can_fold":  false,
			},
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, c.manager.GetPlayerActions())
	}
}

func TestAllInPlayerOptions(t *testing.T) {
	gamePlayers := []*Player{
		createTestPlayer(t, 500),
		createTestPlayer(t, 750),
		createTestPlayer(t, 50),
		createTestPlayer(t, 1000),
	}
	allInManager := createTestManager(t, gamePlayers...)
	assert.Equal(t, expectedPlayerActions(true, true, false, true), allInManager.GetPlayerActions())

	allInManager.PlayerAction(gamePlayers[0].ID, Action{ActionBet, int64(50)})
	assert.Equal(t, expectedPlayerActions(true, true, false, true), allInManager.GetPlayerActions())

	allInManager.PlayerAction(gamePlayers[1].ID, Action{ActionBet, int64(750)})
	assert.Equal(t, expectedPlayerActions(false, true, false, true), allInManager.GetPlayerActions())

	allInManager.PlayerAction(gamePlayers[2].ID, Action{ActionFold, int64(0)})
	assert.Equal(t, expectedPlayerActions(false, true, false, true), allInManager.GetPlayerActions())

	allInManager.PlayerAction(gamePlayers[3].ID, Action{ActionCall, int64(0)})
	assert.Equal(t, gamePlayers[0].ID, allInManager.State.Table.GetActivePlayer().ID)
	assert.Equal(t, expectedPlayerActions(false, true, false, true), allInManager.GetPlayerActions())

	allInManager.PlayerAction(gamePlayers[0].ID, Action{ActionCall, int64(0)})
	assert.Equal(t, expectedPlayerActions(false, false, false, false), allInManager.GetPlayerActions())
	assert.False(t, allInManager.State.Table.GetActivePlayer().HasOptions())
	assert.True(t, allInManager.IsAllIn())
}
