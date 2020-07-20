package holdem

import (
	"qpoker/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestPlayer(t *testing.T, score int) *Player {
	player := models.CreateTestPlayer()

	return &Player{
		ID:       player.ID,
		Username: player.Username,
		Score:    score,
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

func TestPlay4Hands(t *testing.T) {
	var nilMap map[string]bool
	var player *models.Player

	players := []*Player{}
	for i := 0; i < 4; i++ {
		player = models.CreateTestPlayer()
		players = append(players, &Player{ID: player.ID})
	}
	game := models.CreateTestGame(player.ID)

	gm, err := NewGameManager(game.ID, players, models.GameOptions{Capacity: 5, BigBlind: 100})
	assert.NoError(t, err)

	err = gm.NextHand()
	assert.NoError(t, err)

	// check game hand saved expected values
	assert.Greater(t, gm.gameHand.ID, int64(0))
	assert.Equal(t, gm.gameHand.GameID, game.ID)
	for i := 0; i < 4; i++ {
		assert.Greater(t, gm.gamePlayerHands[players[i].ID].ID, int64(0))
		assert.Equal(t, gm.gamePlayerHands[players[i].ID].GameHandID, gm.gameHand.ID)
		assert.Equal(t, players[i].ID, gm.gamePlayerHands[players[i].ID].PlayerID)
		assert.Equal(t, 13, len(gm.gamePlayerHands[players[i].ID].Cards))
		assert.Equal(t, int64(0), gm.gamePlayerHands[players[i].ID].StartingStack)
	}

	// Try to move out of turn
	for _, player := range players {
		assert.Equal(t, nilMap, player.Options)
		_, err = gm.PlayerAction(player.ID, NewActionPlay("AC"))
		assert.Error(t, err)
	}

	// assert.True(t, players[4].Options["can_play"])
	// complete, err := gm.PlayerAction(players[4].ID, NewActionFold())
	// assert.NoError(t, err)
	// assert.False(t, complete)
	// assert.Equal(t, nilMap, players[4].Options)

	// assert.True(t, players[0].Options["can_bet"])
	// assert.True(t, players[0].Options["can_call"])
	// assert.True(t, players[0].Options["can_fold"])
	// assert.False(t, players[0].Options["can_check"])
	// complete, err = gm.PlayerAction(players[0].ID, NewActionFold())
	// assert.NoError(t, err)
	// assert.False(t, complete)
	// assert.Equal(t, nilMap, players[0].Options)

	// complete, err = gm.PlayerAction(players[1].ID, NewActionFold())
	// assert.NoError(t, err)
	// assert.False(t, complete)
	// complete, err = gm.PlayerAction(players[2].ID, NewActionFold())
	// assert.NoError(t, err)
	// assert.True(t, complete)

	// // check game hand final save expected values
	// expectedFinalStacks := []int64{200, 200, 150, 250, 200}
	// assert.Equal(t, 0, len(gm.gameHand.Board))
	// assert.Equal(t, int64(150), gm.gameHand.Payouts[players[3].ID])
	// assert.Equal(t, int64(100), gm.gameHand.Bets[players[3].ID])
	// assert.Equal(t, int64(50), gm.gameHand.Bets[players[2].ID])
	// assert.Equal(t, gm.gameHand.GameID, game.ID)
	// for i := 0; i < 5; i++ {
	// 	assert.Greater(t, gm.gamePlayerHands[players[i].ID].ID, int64(0))
	// 	assert.Equal(t, 2, len(gm.gamePlayerHands[players[i].ID].Cards))
	// 	assert.Equal(t, int64(200), gm.gamePlayerHands[players[i].ID].StartingStack)
	// 	assert.Equal(t, expectedFinalStacks[i], gm.gamePlayerHands[players[i].ID].EndingStack, fmt.Sprintf("Final Player Hand: pos(%d) %+v\n", i, gm.gamePlayerHands[players[i].ID]))
	// 	assert.False(t, players[i].CardsVisible)
	// }
}
