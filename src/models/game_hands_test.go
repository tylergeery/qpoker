package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGameHandBy(t *testing.T) {
	// Given
	player := CreateTestPlayer()
	player2 := CreateTestPlayer()
	game := CreateTestGame(player)
	gameHand := &GameHand{
		GameID: game.ID,
		Board:  []string{"2C", "AD", "KD", "QC", "2D"},
		Payouts: UserStackMap{
			player.ID: 500,
		},
		Bets: UserStackMap{
			player.ID:  250,
			player2.ID: 250,
		},
	}
	err := gameHand.Save()
	assert.NoError(t, err)
	err = gameHand.Save()
	assert.NoError(t, err)

	// When
	fetchedGameHand, err := GetGameHandBy("id", gameHand.ID)
	assert.NoError(t, err)

	// Then
	assert.Equal(t, fetchedGameHand.ID, gameHand.ID)
	assert.Equal(t, fetchedGameHand.GameID, gameHand.GameID)
	assert.Equal(t, fetchedGameHand.Board, gameHand.Board)
	assert.Equal(t, fetchedGameHand.Payouts, gameHand.Payouts)
	assert.Equal(t, fetchedGameHand.Bets, gameHand.Bets)
	assert.Equal(t, fetchedGameHand.CreatedAt.Unix(), gameHand.CreatedAt.Unix())
	assert.Equal(t, fetchedGameHand.UpdatedAt.Unix(), gameHand.UpdatedAt.Unix())
}

func TestGetHandsForGame(t *testing.T) {
	// Given
	player := CreateTestPlayer()
	player2 := CreateTestPlayer()
	game := CreateTestGame(player)
	gameHands := []*GameHand{
		&GameHand{
			GameID: game.ID,
			Board:  []string{"2C", "AD", "KD", "QC", "2D"},
			Payouts: UserStackMap{
				player.ID: 500,
			},
			Bets: UserStackMap{
				player.ID:  250,
				player2.ID: 250,
			},
		},
		&GameHand{
			GameID: game.ID,
			Board:  []string{"2C", "AD", "KD", "QC", "2D"},
			Payouts: UserStackMap{
				player.ID: 250,
			},
			Bets: UserStackMap{
				player.ID: 250,
			},
		},
		&GameHand{
			GameID: game.ID,
			Board:  []string{"TD", "AD", "KD", "4C", "2D"},
			Payouts: UserStackMap{
				player2.ID: 200,
			},
			Bets: UserStackMap{
				player.ID:  100,
				player2.ID: 100,
			},
		},
	}
	for _, gameHand := range gameHands {
		err := gameHand.Save()
		assert.NoError(t, err)
		err = gameHand.Save()
		assert.NoError(t, err)
	}

	// When
	fetchedHands, err := GetHandsForGame(game.ID, int64(0), 2)
	assert.NoError(t, err)
	fetchedHands2, err := GetHandsForGame(game.ID, fetchedHands[1].ID, 2)
	assert.NoError(t, err)
	totalHands := append(fetchedHands, fetchedHands2...)

	// Then
	for i := range totalHands {
		assert.Equal(t, totalHands[i].ID, gameHands[i].ID)
		assert.Equal(t, totalHands[i].GameID, gameHands[i].GameID)
		assert.Equal(t, totalHands[i].Board, gameHands[i].Board)
		assert.Equal(t, totalHands[i].Payouts, gameHands[i].Payouts)
		assert.Equal(t, totalHands[i].Bets, gameHands[i].Bets)
		assert.Equal(t, totalHands[i].CreatedAt.Unix(), gameHands[i].CreatedAt.Unix())
		assert.Equal(t, totalHands[i].UpdatedAt.Unix(), gameHands[i].UpdatedAt.Unix())
	}
}
