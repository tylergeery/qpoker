package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGamePlayerHandsCrud(t *testing.T) {
	player1 := CreateTestPlayer()
	player2 := CreateTestPlayer()
	player3 := CreateTestPlayer()
	game := CreateTestGame(player1.ID)
	hand := &GameHand{GameID: game.ID}
	err := hand.Save()
	assert.NoError(t, err)

	fetchedHand, err := GetGameHandBy("id", hand.ID)
	assert.NoError(t, err)
	assert.Greater(t, hand.ID, int64(0))
	assert.Equal(t, hand.ID, fetchedHand.ID)
	assert.Equal(t, game.ID, fetchedHand.GameID)
	assert.Equal(t, hand.CreatedAt, fetchedHand.CreatedAt)
	assert.Equal(t, hand.UpdatedAt, fetchedHand.UpdatedAt)

	playerHand1 := &GamePlayerHand{
		GameHandID:    hand.ID,
		PlayerID:      player1.ID,
		StartingStack: int64(50),
	}
	err = playerHand1.Save()
	assert.NoError(t, err)
	playerHand2 := &GamePlayerHand{
		GameHandID:    hand.ID,
		PlayerID:      player2.ID,
		StartingStack: int64(150),
	}
	err = playerHand2.Save()
	assert.NoError(t, err)
	playerHand3 := &GamePlayerHand{
		GameHandID:    hand.ID,
		PlayerID:      player3.ID,
		StartingStack: int64(150),
	}
	err = playerHand3.Save()
	assert.NoError(t, err)

	fetchedPlayerHand1, err := GetGamePlayerHandBy("id", playerHand1.ID)
	assert.NoError(t, err)
	assert.Equal(t, playerHand1.ID, fetchedPlayerHand1.ID)
	assert.Equal(t, playerHand1.GameHandID, fetchedPlayerHand1.GameHandID)
	assert.Equal(t, playerHand1.PlayerID, fetchedPlayerHand1.PlayerID)
	assert.Equal(t, playerHand1.StartingStack, fetchedPlayerHand1.StartingStack)

	fetchedPlayerHand2, err := GetGamePlayerHandBy("id", playerHand2.ID)
	assert.NoError(t, err)
	assert.Equal(t, playerHand2.ID, fetchedPlayerHand2.ID)
	assert.Equal(t, playerHand2.GameHandID, fetchedPlayerHand2.GameHandID)
	assert.Equal(t, playerHand2.PlayerID, fetchedPlayerHand2.PlayerID)
	assert.Equal(t, playerHand2.StartingStack, fetchedPlayerHand2.StartingStack)

	fetchedPlayerHand3, err := GetGamePlayerHandBy("id", playerHand3.ID)
	assert.NoError(t, err)
	assert.Equal(t, playerHand3.ID, fetchedPlayerHand3.ID)
	assert.Equal(t, playerHand3.GameHandID, fetchedPlayerHand3.GameHandID)
	assert.Equal(t, playerHand3.PlayerID, fetchedPlayerHand3.PlayerID)
	assert.Equal(t, playerHand3.StartingStack, fetchedPlayerHand3.StartingStack)
	assert.Equal(t, fetchedPlayerHand3.CreatedAt.Unix(), playerHand3.CreatedAt.Unix())
	assert.Equal(t, fetchedPlayerHand3.UpdatedAt.Unix(), playerHand3.UpdatedAt.Unix())

	// Update hand on completion
	playerHand3.Cards = []string{"1A", "5D"}
	playerHand3.EndingStack = int64(25)
	err = playerHand3.Save()
	assert.NoError(t, err)
	fetchedPlayerHand3, err = GetGamePlayerHandBy("id", playerHand3.ID)
	assert.NoError(t, err)
	assert.Equal(t, playerHand3.Cards, fetchedPlayerHand3.Cards)
	assert.Equal(t, playerHand3.EndingStack, fetchedPlayerHand3.EndingStack)
	assert.Equal(t, fetchedPlayerHand3.UpdatedAt.Unix(), playerHand3.UpdatedAt.Unix())
}
