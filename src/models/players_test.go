package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPlayerFromID(t *testing.T) {
	// Given
	player := &Player{
		Username: fmt.Sprintf("Test %d", time.Now().UTC().UnixNano()),
		Email:    fmt.Sprintf("token%d", time.Now().UTC().UnixNano()),
	}
	player.Create("testpw")

	// When
	fetchedPlayer, _ := GetPlayerFromID(player.ID)

	// Then
	assert.Equal(t, fetchedPlayer.ID, player.ID, fmt.Sprintf("Expected fetchedPlayer.ID (%d) to match player.ID (%d)", fetchedPlayer.ID, player.ID))
	assert.Equal(t, fetchedPlayer.Username, player.Username, fmt.Sprintf("Expected fetchedPlayer.Username (%s) to match player.Username (%s)", fetchedPlayer.Username, player.Username))
	assert.Equal(t, fetchedPlayer.Email, player.Email, fmt.Sprintf("Expected fetchedPlayer.Email (%s) to match player.Username (%s)", fetchedPlayer.Email, player.Email))
	assert.Equal(t, fetchedPlayer.pw, "", fmt.Sprintf("Expected fetchedPlayer.pw (%s) to be blank", fetchedPlayer.Email))
	assert.Equal(t, fetchedPlayer.CreatedAt.Unix(), player.CreatedAt.Unix(), fmt.Sprintf("Expected fetchedPlayer.CreatedAt (%d) to match player.CreatedAt (%d)", fetchedPlayer.CreatedAt.Unix(), player.CreatedAt.Unix()))
	assert.Equal(t, fetchedPlayer.UpdatedAT.Unix(), player.UpdatedAt.Unix(), fmt.Sprintf("Expected fetchedPlayer.UpdatedAt (%d) to match player.UpdatedAt (%d)", fetchedPlayer.UpdatedAt.Unix(), player.UpdatedAt.Unix()))
}
