package models

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetGameBy(t *testing.T) {
	// Given
	rand.Seed(time.Now().UTC().UnixNano())
	pw := "asdfja193242"
	player := &Player{
		Username: fmt.Sprintf("Test Player %d", time.Now().UTC().UnixNano()),
		Email:    fmt.Sprintf("token%d", time.Now().UTC().UnixNano()),
	}
	player.Create(pw)
	game := &Game{
		Name:     fmt.Sprintf("Testing Game %d", time.Now().UTC().UnixNano()),
		Capacity: 12,
		OwnerID:  player.ID,
	}
	game.Save()

	// When
	fetchedGame, _ := GetGameBy("id", game.ID)

	// Then
	assert.Equal(t, fetchedGame.ID, game.ID)
	assert.Equal(t, fetchedGame.Name, game.Name)
	assert.Equal(t, fetchedGame.Slug, game.Slug)
	assert.True(t, len(fetchedGame.Slug) == 16)
	assert.Equal(t, fetchedGame.Capacity, game.Capacity)
	assert.Equal(t, fetchedGame.CreatedAt.Unix(), game.CreatedAt.Unix())
	assert.Equal(t, fetchedGame.UpdatedAt.Unix(), game.UpdatedAt.Unix())
}
