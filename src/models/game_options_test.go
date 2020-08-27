package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameOptionsCrud(t *testing.T) {
	// Given
	player := CreateTestPlayer()
	game := CreateTestGame(player.ID, 1)
	record := GameOptions{
		GameID:     game.ID,
		GameTypeID: game.GameTypeID,
		Options: map[string]interface{}{
			"capacity":  int64(5),
			"big_blind": int64(1000),
		},
	}
	err := record.Save()
	assert.NoError(t, err)

	// When
	options, err := GetGameOptionsForGame(game.ID, game.GameTypeID)
	assert.NoError(t, err)

	// Then
	assert.Equal(t, options.GameID, game.ID)
	assert.Equal(t, options.Options["capacity"], record.Options["capacity"])
	assert.Equal(t, options.Options["big_blind"], record.Options["big_blind"])
	assert.Equal(t, int64(5), options.Options["capacity"])
	assert.Equal(t, int64(1000), options.Options["big_blind"])

	// Update Options
	record.Options["capacity"] = int64(10)
	record.Options["big_blind"] = int64(200)
	err = record.Save()
	assert.NoError(t, err)

	options, err = GetGameOptionsForGame(game.ID, game.GameTypeID)
	assert.NoError(t, err)

	// Then
	assert.Equal(t, options.GameID, game.ID)
	assert.Equal(t, options.Options["capacity"], record.Options["capacity"])
	assert.Equal(t, options.Options["big_blind"], record.Options["big_blind"])
	assert.Equal(t, int64(10), options.Options["capacity"])
	assert.Equal(t, int64(200), options.Options["big_blind"])
}
