package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGameBy(t *testing.T) {
	// Given
	player := CreateTestPlayer()
	game := CreateTestGame(player.ID, 2)

	// When
	fetchedGame, err := GetGameBy("id", game.ID)
	assert.NoError(t, err)

	// Then
	assert.Equal(t, fetchedGame.ID, game.ID)
	assert.Equal(t, fetchedGame.Name, game.Name)
	assert.Equal(t, fetchedGame.Slug, game.Slug)
	assert.True(t, len(fetchedGame.Slug) >= 16)
	assert.Equal(t, fetchedGame.Status, game.Status)
	assert.Equal(t, fetchedGame.CreatedAt.Unix(), game.CreatedAt.Unix())
	assert.Equal(t, fetchedGame.UpdatedAt.Unix(), game.UpdatedAt.Unix())

	// When
	fetchedGame, err = GetGameBy("slug", game.Slug)
	assert.NoError(t, err)

	// Then
	assert.Equal(t, fetchedGame.ID, game.ID)
	assert.Equal(t, fetchedGame.Name, game.Name)
	assert.Equal(t, fetchedGame.Slug, game.Slug)
	assert.True(t, len(fetchedGame.Slug) >= 16)
	assert.Equal(t, fetchedGame.CreatedAt.Unix(), game.CreatedAt.Unix())
	assert.Equal(t, fetchedGame.UpdatedAt.Unix(), game.UpdatedAt.Unix())

	// Update
	oldGameName := game.Name
	game.Name = fmt.Sprintf("Updated %s", oldGameName)
	game.Save()

	fetchedGame, _ = GetGameBy("name", game.Name)
	_, err = GetGameBy("name", oldGameName)

	// Then
	assert.Error(t, err)
	assert.Equal(t, fetchedGame.ID, game.ID)
	assert.Equal(t, fetchedGame.Name, game.Name)
	assert.Equal(t, fetchedGame.Slug, game.Slug)
	assert.True(t, len(fetchedGame.Slug) >= 16)
	assert.Equal(t, fetchedGame.CreatedAt.Unix(), game.CreatedAt.Unix())
	assert.Equal(t, fetchedGame.UpdatedAt.Unix(), game.UpdatedAt.Unix())
}
