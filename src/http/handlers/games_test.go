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

func i64(value interface{}) int64 {
	v := value.(float64)

	return int64(v)
}

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
		"name":         name,
		"game_type_id": 1,
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
	response, err := app.Test(req, -1)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &game)

	assert.Equal(t, 201, response.StatusCode)
	assert.Greater(t, game.ID, int64(0))
	assert.Equal(t, name, game.Name)
	assert.True(t, len(game.Slug) >= 15)
	assert.Equal(t, game.Status, models.GameStatusInit)
	assert.Equal(t, float64(capacity), game.Options["capacity"])
	assert.Equal(t, float64(bigBlind), game.Options["big_blind"])
	assert.Equal(t, float64(5), game.Options["time_between_hands"])
	assert.Equal(t, float64(5000), game.Options["buy_in_max"])
	assert.Equal(t, float64(500), game.Options["buy_in_min"])
	assert.Greater(t, game.CreatedAt.Unix(), int64(0))
	assert.Greater(t, game.UpdatedAt.Unix(), int64(0))
}

func TestGameUpdateSuccess(t *testing.T) {
	// Given
	var updatedGame models.Game

	player := test.CreateTestPlayer()
	game := test.CreateTestGame(player, 1)
	body := map[string]interface{}{
		"name": "temp " + game.Name,
		"options": map[string]interface{}{
			"capacity":   10,
			"big_blind":  60,
			"buy_in_min": 5,
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
	assert.Equal(t, float64(10), updatedGame.Options["capacity"])
	assert.Equal(t, float64(60), updatedGame.Options["big_blind"])
	assert.Equal(t, float64(5), updatedGame.Options["buy_in_min"])
	assert.Equal(t, game.Slug, updatedGame.Slug)
	assert.Equal(t, updatedGame.CreatedAt.Unix(), game.CreatedAt.Unix())
	assert.GreaterOrEqual(t, updatedGame.UpdatedAt.Unix(), game.UpdatedAt.Unix())
}

func TestGetGameFailure(t *testing.T) {
	player := test.CreateTestPlayer()
	player2 := test.CreateTestPlayer()
	game := test.CreateTestGame(player, 1)

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
	game := test.CreateTestGame(player, 1)
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
	assert.Equal(t, retrievedGame.Options["capacity"], float64(game.Options["capacity"].(int64)))
	assert.Equal(t, game.CreatedAt.Unix(), retrievedGame.CreatedAt.Unix())
	assert.Equal(t, game.UpdatedAt.Unix(), retrievedGame.CreatedAt.Unix())
}

func TestGetGameHistoryEmpty(t *testing.T) {
	// Given
	var gameHistory []interface{}

	player := test.CreateTestPlayer()
	game := test.CreateTestGame(player, 1)
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
	game := test.CreateTestGame(player, 1)
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
		PlayerID: player2.ID,
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
				GameHandID: hand.ID,
				PlayerID:   player.ID,
				Cards:      []string{"JS", "JC"},
				Starting:   int64(120 + (i * 20)),
				Ending:     int64(120 + ((i + 1) * 20)),
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
	chipRequest := gameHistory[0].(map[string]interface{})
	assert.Equal(t, int64(120), i64(chipRequest["amount"]))
	assert.Equal(t, game.ID, i64(chipRequest["game_id"]))
	assert.Equal(t, player.ID, i64(chipRequest["player_id"]))
	assert.Equal(t, "approved", chipRequest["status"].(string))

	chipRequest = gameHistory[1].(map[string]interface{})
	assert.Equal(t, int64(360), i64(chipRequest["amount"]))
	assert.Equal(t, game.ID, i64(chipRequest["game_id"]))
	assert.Equal(t, player2.ID, i64(chipRequest["player_id"]))
	assert.Equal(t, "approved", chipRequest["status"].(string))

	hand := gameHistory[2].(map[string]interface{})
	assert.Equal(t, nil, hand["board"])
	assert.Equal(t, []interface{}{"JS", "JC"}, hand["cards"])
	assert.Equal(t, map[string]interface{}{"Int64": float64(140), "Valid": true}, hand["ending"])
	assert.Equal(t, map[string]interface{}{"Int64": float64(120), "Valid": true}, hand["starting"])
	assert.Equal(t, game.ID, i64(hand["game_id"]))

	hand = gameHistory[3].(map[string]interface{})
	assert.Equal(t, nil, hand["board"])
	assert.Equal(t, nil, hand["cards"])
	assert.Equal(t, map[string]interface{}{"Int64": float64(0), "Valid": false}, hand["ending"])
	assert.Equal(t, map[string]interface{}{"Int64": float64(0), "Valid": false}, hand["starting"])
	assert.Equal(t, game.ID, i64(hand["game_id"]))

	hand = gameHistory[4].(map[string]interface{})
	assert.Equal(t, nil, hand["board"])
	assert.Equal(t, nil, hand["cards"])
	assert.Equal(t, map[string]interface{}{"Int64": float64(0), "Valid": false}, hand["ending"])
	assert.Equal(t, map[string]interface{}{"Int64": float64(0), "Valid": false}, hand["starting"])
	assert.Equal(t, game.ID, i64(hand["game_id"]))

	chipRequest = gameHistory[5].(map[string]interface{})
	assert.Equal(t, int64(120), i64(chipRequest["amount"]))
	assert.Equal(t, game.ID, i64(chipRequest["game_id"]))
	assert.Equal(t, player.ID, i64(chipRequest["player_id"]))
	assert.Equal(t, "", chipRequest["status"].(string))
}
