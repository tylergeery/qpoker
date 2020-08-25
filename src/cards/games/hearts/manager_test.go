package hearts

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
		Score:    int64(score),
	}
}

func createTestManager(t *testing.T, players ...*Player) *GameManager {
	game := models.CreateTestGame(players[0].ID)

	manager, err := NewGameManager(game.ID, players, GameOptions{})
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

	gm, err := NewGameManager(game.ID, players, GameOptions{Capacity: 4})
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
		assert.Equal(t, int64(0), gm.gamePlayerHands[players[i].ID].Starting)
	}

	// Try to move out of turn
	for _, player := range players {
		assert.Equal(t, nilMap, player.Options)
		_, err = gm.PlayerAction(player.ID, NewActionPlay(player.Cards[0].ToString()))
		assert.Error(t, err)
	}

	for j := 0; j < 4; j++ {
		// Pass cards
		if j != 3 {
			assert.Equal(t, gm.Status, StatusReady)
			for i := range players {
				player := players[(i+j)%4] // Let arbitrary player pass first
				err = gm.PlayerPass(player.ID, player.Cards[:3])
				assert.NoError(t, err)
			}
		}
		assert.Equal(t, gm.Status, StatusActive)

		// Play round
		for i := 0; i < 52; i++ {
			player := gm.State.Table.GetActivePlayer()
			_, err = gm.PlayerAction(player.ID, NewActionPlay(player.Cards[0].ToString()))
			assert.NoError(t, err)
		}
		assert.True(t, gm.isComplete())

		// Count save hearts (ensure all accounted for)
		totalPoints := int64(0)
		for _, hand := range gm.gamePlayerHands {
			totalPoints += hand.Ending - hand.Starting
		}

		assert.True(t, totalPoints == int64(26) || totalPoints == int64(78))

		// Proceed to next hand
		err = gm.NextHand()
		assert.NoError(t, err)
	}
}
