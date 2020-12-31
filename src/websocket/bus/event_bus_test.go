package bus

import (
	"qpoker/cards/games/holdem"
	"qpoker/models"
	"qpoker/websocket/connection"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReloadGameState(t *testing.T) {
	player := models.CreateTestPlayer()
	game := models.CreateTestGame(player.ID, 1)
	game.Status = models.GameStatusStart
	game.Save()
	client := &connection.Client{PlayerID: player.ID, GameID: game.ID}
	eventBus := NewEventBus()

	err := eventBus.loadGameState(client)
	loadedGame, ok := eventBus.games[game.ID]

	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, models.GameStatusStart, loadedGame.Data().Game.Status)
}

func TestReloadPlayerStack(t *testing.T) {
	player := models.CreateTestPlayer()
	game := models.CreateTestGame(player.ID, 1)
	gameChipRequest := &models.GameChipRequest{
		GameID:   game.ID,
		PlayerID: player.ID,
		Amount:   int64(50),
		Status:   models.GameChipRequestStatusApproved,
	}
	gameChipRequest.Save()
	gameHand := &models.GameHand{
		GameID: game.ID,
		Board:  []string{"JC", "JD"},
	}
	gameHand.Save()
	gamePlayerHand := &models.GamePlayerHand{
		GameHandID: gameHand.ID,
		PlayerID:   player.ID,
		Starting:   int64(100),
		Ending:     int64(25),
	}
	gamePlayerHand.Save()
	gameChipRequest = &models.GameChipRequest{
		GameID:   game.ID,
		PlayerID: player.ID,
		Amount:   int64(100),
		Status:   models.GameChipRequestStatusApproved,
	}
	gameChipRequest.Save()

	client := &connection.Client{PlayerID: player.ID, GameID: game.ID}
	eventBus := NewEventBus()
	eventBus.loadGameState(client)

	controller := eventBus.games[game.ID]
	gamePlayer := controller.AddPlayer(player).(*holdem.Player)

	assert.Equal(t, int64(125), gamePlayer.Stack)
}
