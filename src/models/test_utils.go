package models

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var persistentSeedOnce sync.Once

func setSeed() {
	persistentSeedOnce.Do(func() {
		rand.Seed(time.Now().UTC().UnixNano())
	})
}

// CreateTestPlayer creates a player for a test case
func CreateTestPlayer() *Player {
	setSeed()

	pw := "asdfja193242"
	player := &Player{
		Username: fmt.Sprintf("Test Player %d", time.Now().UTC().UnixNano()),
		Email:    fmt.Sprintf("token%d", time.Now().UTC().UnixNano()),
	}
	player.Create(pw)

	return player
}

// CreateTestGame creates a game for a test game
func CreateTestGame(playerID int64) *Game {
	setSeed()

	game := &Game{
		Name:    fmt.Sprintf("Testing Game %d", time.Now().UTC().UnixNano()),
		OwnerID: playerID,
		Status:  GameStatusInit,
	}
	game.Save()

	return game
}
