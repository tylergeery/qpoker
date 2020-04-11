package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameOptionsCrud(t *testing.T) {
	// Given
	player := CreateTestPlayer()
	game := CreateTestGame(player)
	record := GameOptionsRecord{
		GameID: game.ID,
		Options: GameOptions{
			Capacity: 5,
			BigBlind: int64(1000),
		},
	}
	err := record.Save()
	assert.NoError(t, err)

	// When
	fetchedRecord, err := GetGameOptionsRecordBy("id", record.ID)
	assert.NoError(t, err)
	options, err := GetGameOptionsForGame(game.ID)
	assert.NoError(t, err)

	// Then
	assert.Greater(t, fetchedRecord.ID, int64(0))
	assert.Equal(t, fetchedRecord.ID, record.ID)
	assert.Equal(t, fetchedRecord.GameID, game.ID)
	assert.Equal(t, fetchedRecord.Options, record.Options)
	assert.Equal(t, fetchedRecord.CreatedAt.Unix(), record.CreatedAt.Unix())
	assert.Equal(t, fetchedRecord.UpdatedAt.Unix(), record.UpdatedAt.Unix())
	assert.Equal(t, record.Options, options)
	assert.Equal(t, 5, options.Capacity)
	assert.Equal(t, int64(1000), options.BigBlind)

	// Update Options
	record.Options.Capacity = 10
	record.Options.BigBlind = int64(200)
	err = record.Save()
	assert.NoError(t, err)

	fetchedRecord, err = GetGameOptionsRecordBy("id", record.ID)
	assert.NoError(t, err)
	options, err = GetGameOptionsForGame(game.ID)
	assert.NoError(t, err)

	// Then
	assert.Greater(t, fetchedRecord.ID, int64(0))
	assert.Equal(t, fetchedRecord.ID, record.ID)
	assert.Equal(t, fetchedRecord.GameID, game.ID)
	assert.Equal(t, fetchedRecord.Options, record.Options)
	assert.Equal(t, fetchedRecord.CreatedAt.Unix(), record.CreatedAt.Unix())
	assert.Equal(t, fetchedRecord.UpdatedAt.Unix(), record.UpdatedAt.Unix())
	assert.Equal(t, record.Options, options)
	assert.Equal(t, 10, options.Capacity)
	assert.Equal(t, int64(200), options.BigBlind)
}
