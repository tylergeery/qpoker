package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"qpoker/models"
	"qpoker/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGameCreateInvalid(t *testing.T) {
	type TestCase struct {
		body     map[string]interface{}
		headers  map[string]string
		expected int
	}

	player := test.CreateTestPlayer()
	cases := []TestCase{
		TestCase{
			body: map[string]interface{}{
				"name":     "Game test",
				"capacity": 4,
			},
			headers:  map[string]string{},
			expected: 403,
		},
		TestCase{
			body: map[string]interface{}{},
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", player.Token),
			},
			expected: 400,
		},
		TestCase{
			body: map[string]interface{}{
				"name":     "Game test",
				"capacity": -2,
			},
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", player.Token),
			},
			expected: 400,
		},
		TestCase{
			body: map[string]interface{}{
				"name":     "small",
				"capacity": 8,
			},
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", player.Token),
			},
			expected: 400,
		},
	}

	app := CreateApp()

	for _, c := range cases {
		req := test.CreateTestRequest("POST", "/api/v1/games", c.headers, c.body)

		response, err := app.Test(req)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(response.Body)
		assert.NoError(t, err)

		assert.Equal(t, c.expected, response.StatusCode)

		var responseMap map[string][]string
		json.Unmarshal(body, &responseMap)

		assert.NotEqual(t, "", responseMap["errors"][0])
	}
}

func TestCreateGameSuccess(t *testing.T) {
	// Given
	var game models.Game
	player := test.CreateTestPlayer()
	ts := time.Now().UTC().UnixNano()

	name := fmt.Sprintf("Test Game %d", ts)
	capacity := 20
	bigBlind := int64(100)
	body := map[string]interface{}{
		"name": name,
		"options": map[string]interface{}{
			"capacity":  capacity,
			"big_blind": bigBlind,
		},
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", player.Token),
	}
	req := test.CreateTestRequest("POST", "/api/v1/games", headers, body)
	app := CreateApp()

	// When
	response, err := app.Test(req)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &game)

	assert.Equal(t, 201, response.StatusCode)
	assert.Greater(t, game.ID, int64(0))
	assert.Equal(t, name, game.Name)
	assert.True(t, len(game.Slug) >= 16)
	assert.Equal(t, game.Status, models.GameStatusInit)
	assert.Equal(t, capacity, game.Options.Capacity)
	assert.Equal(t, bigBlind, game.Options.BigBlind)
	assert.Equal(t, 5, game.Options.TimeBetweenHands)
	assert.Equal(t, int64(10)*bigBlind, game.Options.BuyInMax)
	assert.Equal(t, bigBlind, game.Options.BuyInMin)
	assert.Greater(t, game.CreatedAt.Unix(), int64(0))
	assert.Greater(t, game.UpdatedAt.Unix(), int64(0))
}

func TestGameUpdateSuccess(t *testing.T) {
	// Given
	var updatedGame models.Game

	player := test.CreateTestPlayer()
	game := test.CreateTestGame(player)
	body := map[string]interface{}{
		"name": "temp " + game.Name,
		"options": map[string]interface{}{
			"capacity":  10,
			"big_blind": 60,
		},
	}
	req := test.CreateTestRequest("PUT", fmt.Sprintf("/api/v1/games/%d", game.ID), map[string]string{"Authorization": fmt.Sprintf("Bearer %s", player.Token)}, body)
	app := CreateApp()

	// When
	response, err := app.Test(req)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &updatedGame)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, game.ID, updatedGame.ID)
	assert.Equal(t, body["name"], updatedGame.Name)
	assert.Equal(t, 10, updatedGame.Options.Capacity)
	assert.Equal(t, int64(60), updatedGame.Options.BigBlind)
	assert.Equal(t, game.Slug, updatedGame.Slug)
	assert.Equal(t, updatedGame.CreatedAt.Unix(), game.CreatedAt.Unix())
	assert.GreaterOrEqual(t, updatedGame.UpdatedAt.Unix(), game.UpdatedAt.Unix())
}

func TestGetGameFailure(t *testing.T) {
	player := test.CreateTestPlayer()
	player2 := test.CreateTestPlayer()
	game := test.CreateTestGame(player)

	type TestCase struct {
		headers  map[string]string
		expected int
	}
	cases := []TestCase{
		TestCase{
			headers:  map[string]string{},
			expected: 403,
		},
		TestCase{
			headers:  map[string]string{"Content-Type": "application/json", "Authorization": "Faketoken.faker"},
			expected: 403,
		},
		TestCase{
			headers:  map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", player2.Token)},
			expected: 404,
		},
	}
	app := CreateApp()

	for _, c := range cases {
		// When
		req := test.CreateTestRequest("GET", fmt.Sprintf("/api/v1/games/%d", game.ID), c.headers, nil)
		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, c.expected, response.StatusCode, c.headers["Authorization"])
	}
}

func TestGetGameSuccess(t *testing.T) {
	// Given
	var retrievedGame models.Game

	player := test.CreateTestPlayer()
	game := test.CreateTestGame(player)
	headers := map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", player.Token)}
	req := test.CreateTestRequest("GET", fmt.Sprintf("/api/v1/games/%d", game.ID), headers, nil)
	app := CreateApp()

	// When
	response, err := app.Test(req)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &retrievedGame)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, retrievedGame.ID, game.ID)
	assert.Equal(t, retrievedGame.Name, game.Name)
	assert.Equal(t, retrievedGame.Slug, game.Slug)
	assert.Equal(t, retrievedGame.Options.Capacity, game.Options.Capacity)
	assert.Equal(t, game.CreatedAt.Unix(), retrievedGame.CreatedAt.Unix())
	assert.Equal(t, game.UpdatedAt.Unix(), retrievedGame.CreatedAt.Unix())
}

func TestGetGameHistoryEmpty(t *testing.T) {
	// Given
	var gameHistory []interface{}

	player := test.CreateTestPlayer()
	game := test.CreateTestGame(player)
	headers := map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", player.Token)}
	req := test.CreateTestRequest("GET", fmt.Sprintf("/api/v1/games/%d/history", game.ID), headers, nil)
	app := CreateApp()

	// When
	response, err := app.Test(req)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &gameHistory)
	assert.Equal(t, 0, len(gameHistory))
}

func TestGetGameHistorySuccess(t *testing.T) {
	// Given
	var gameHistory []interface{}

	// Create GameHistory
	player := test.CreateTestPlayer()
	player2 := test.CreateTestPlayer()
	game := test.CreateTestGame(player)
	req1 := &models.GameChipRequest{
		GameID:   game.ID,
		PlayerID: player.ID,
		Amount:   120,
		Status:   models.GameChipRequestStatusApproved,
	}
	err := req1.Save()
	assert.NoError(t, err)
	req2 := &models.GameChipRequest{
		GameID:   game.ID,
		PlayerID: player.ID,
		Amount:   360,
		Status:   models.GameChipRequestStatusApproved,
	}
	err = req2.Save()
	assert.NoError(t, err)

	for i := 0; i < 3; i++ {
		hand := &models.GameHand{
			GameID: game.ID,
			Board:  []string{"2D", "3C", "4S"},
			Payouts: models.UserStackMap{
				player.ID: 40,
			},
			Bets: models.UserStackMap{
				player2.ID: 20,
				player.ID:  20,
			},
		}
		err = hand.Save()
		assert.NoError(t, err)

		if i == 0 {
			gamePlayerHand := &models.GamePlayerHand{
				GameHandID:    hand.ID,
				PlayerID:      player.ID,
				Cards:         []string{"JS", "JC"},
				StartingStack: int64(120 + (i * 20)),
				EndingStack:   int64(120 + ((i + 1) * 20)),
			}
			err = gamePlayerHand.Save()
			assert.NoError(t, err)
		}
	}
	req3 := &models.GameChipRequest{
		GameID:   game.ID,
		PlayerID: player.ID,
		Amount:   120,
	}
	err = req3.Save()
	assert.NoError(t, err)

	headers := map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", player.Token)}
	req := test.CreateTestRequest("GET", fmt.Sprintf("/api/v1/games/%d/history", game.ID), headers, nil)
	app := CreateApp()

	// When
	response, err := app.Test(req)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &gameHistory)
	assert.Equal(t, 6, len(gameHistory))

	// Assert things about the 6
}
