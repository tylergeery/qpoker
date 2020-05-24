package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGameChipRequestBy(t *testing.T) {
	// Given
	player := CreateTestPlayer()
	player2 := CreateTestPlayer()
	game := CreateTestGame(player.ID)
	req := &GameChipRequest{
		GameID:   game.ID,
		PlayerID: player2.ID,
		Amount:   150,
		Status:   GameChipRequestStatusInit,
	}
	err := req.Save()
	assert.NoError(t, err)

	// When
	fetchedRequest, err := GetGameChipRequestBy("id", req.ID)
	assert.NoError(t, err)

	// Then
	assert.Equal(t, fetchedRequest.ID, req.ID)
	assert.Equal(t, fetchedRequest.PlayerID, req.PlayerID)
	assert.Equal(t, fetchedRequest.PlayerID, req.PlayerID)
	assert.Equal(t, fetchedRequest.Amount, req.Amount)
	assert.Equal(t, fetchedRequest.Status, req.Status)
	assert.Equal(t, fetchedRequest.CreatedAt.Unix(), req.CreatedAt.Unix())
	assert.Equal(t, fetchedRequest.UpdatedAt.Unix(), req.UpdatedAt.Unix())
}
