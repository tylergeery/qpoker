package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"qpoker/auth"
	"qpoker/http/test"
	"qpoker/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlayerCreateInvalid(t *testing.T) {
	type TestCase struct {
		body     map[string]string
		expected string
	}
	cases := []TestCase{
		TestCase{
			body: map[string]string{
				"email": "test@qpoker.com",
				"pw":    "testpassword",
			},
			expected: "{\"errors\": [\"username requires at least 6 characters\"]}",
		},
		TestCase{
			body: map[string]string{
				"username": "testqpokerguy",
				"pw":       "testpassword",
			},
			expected: "{\"errors\": [\"Invalid email format: \"]}",
		},
		TestCase{
			body: map[string]string{
				"email":    "akslfadfalasdfas",
				"username": "testqpokerguy",
				"pw":       "testpassword",
			},
			expected: "{\"errors\": [\"Invalid email format: akslfadfalasdfas\"]}",
		},
		TestCase{
			body: map[string]string{
				"email":    "test@qpoker.com",
				"username": "testqpokerguy",
			},
			expected: "{\"errors\": [\"password requires at least 6 characters\"]}",
		},
	}

	server := CreateApp()

	for _, c := range cases {
		req := test.CreateTestRequest("POST", "/api/v1/players", map[string]string{}, c.body)

		response, err := server.Test(req)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(response.Body)
		assert.NoError(t, err)

		assert.Equal(t, 400, response.StatusCode)
		assert.JSONEq(t, c.expected, string(body))
	}
}

func TestCreatePlayerSuccess(t *testing.T) {
	// Given
	var player models.Player
	ts := time.Now().UTC().UnixNano()
	email := fmt.Sprintf("test+%d@qpoker.com", ts)
	username := fmt.Sprintf("testqpokerguy_%d", ts)
	body := map[string]string{
		"email":    email,
		"username": username,
		"pw":       "testpasstest",
	}
	req := test.CreateTestRequest("POST", "/api/v1/players", nil, body)
	app := CreateApp()

	// When
	response, err := app.Test(req)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &player)
	playerIDFromClaims, _ := auth.GetPlayerIDFromAccessToken(player.Token)

	assert.Equal(t, 201, response.StatusCode)
	assert.Greater(t, player.ID, int64(0))
	assert.Equal(t, email, player.Email)
	assert.Equal(t, username, player.Username)
	assert.Equal(t, player.ID, playerIDFromClaims)
	assert.Greater(t, player.CreatedAt.Unix(), int64(0))
	assert.Greater(t, player.UpdatedAt.Unix(), int64(0))
}

func TestPlayerUpdateSuccess(t *testing.T) {
	// Given
	var updatedPlayer models.Player
	var responseBody map[string]string

	player := test.CreateTestPlayer()
	body := map[string]string{
		"email":    "sa" + player.Email,
		"username": player.Username + "aa",
		"pw":       "newlesspw",
	}
	req := test.CreateTestRequest("PUT", fmt.Sprintf("/api/v1/players/%d", player.ID), map[string]string{"Authorization": fmt.Sprintf("Bearer %s", player.Token)}, body)
	app := CreateApp()

	// When
	response, err := app.Test(req)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &updatedPlayer)
	json.Unmarshal(content, &responseBody)
	playerIDFromClaims, _ := auth.GetPlayerIDFromAccessToken(updatedPlayer.Token)
	_, ok := responseBody["pw"]

	assert.False(t, ok)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, player.ID, updatedPlayer.ID)
	assert.Equal(t, body["email"], updatedPlayer.Email)
	assert.Equal(t, body["username"], updatedPlayer.Username)
	assert.Equal(t, updatedPlayer.ID, playerIDFromClaims)
	assert.Equal(t, updatedPlayer.CreatedAt.Unix(), player.CreatedAt.Unix())
	assert.GreaterOrEqual(t, updatedPlayer.UpdatedAt.Unix(), player.UpdatedAt.Unix())
}

func TestPlayerLoginSuccess(t *testing.T) {
	// Given
	var loggedIn models.Player
	var responseBody map[string]string

	player := test.CreateTestPlayer()
	body := map[string]string{
		"email": player.Email,
		"pw":    test.TestPass,
	}
	req := test.CreateTestRequest("POST", "/api/v1/players/login", nil, body)
	app := CreateApp()

	// When
	response, err := app.Test(req)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &loggedIn)
	json.Unmarshal(content, &responseBody)
	playerIDFromClaims, _ := auth.GetPlayerIDFromAccessToken(loggedIn.Token)
	_, ok := responseBody["pw"]

	assert.False(t, ok)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, loggedIn.ID, player.ID)
	assert.Equal(t, loggedIn.Email, player.Email)
	assert.Equal(t, loggedIn.Username, player.Username)
	assert.Equal(t, player.ID, playerIDFromClaims)
	assert.Equal(t, player.CreatedAt.Unix(), loggedIn.CreatedAt.Unix())
	assert.Equal(t, player.UpdatedAt.Unix(), loggedIn.CreatedAt.Unix())
}

func TestGetPlayerFailure(t *testing.T) {
	player := test.CreateTestPlayer()
	headerAttempts := []map[string]string{
		map[string]string{},
		map[string]string{"Content-Type": "application/json", "Authorization": "Faketoken.faker"},
		map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer a%st", player.Token)},
	}
	app := CreateApp()

	for _, headers := range headerAttempts {
		// When
		req := test.CreateTestRequest("GET", fmt.Sprintf("/api/v1/players/%d", player.ID), headers, nil)
		response, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, response.StatusCode)
	}
}

func TestGetPlayerSuccess(t *testing.T) {
	// Given
	var loggedIn models.Player
	var responseBody map[string]string

	player := test.CreateTestPlayer()
	headers := map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", player.Token)}
	req := test.CreateTestRequest("GET", fmt.Sprintf("/api/v1/players/%d", player.ID), headers, nil)
	app := CreateApp()

	// When
	response, err := app.Test(req)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &loggedIn)
	json.Unmarshal(content, &responseBody)
	playerIDFromClaims, _ := auth.GetPlayerIDFromAccessToken(loggedIn.Token)
	_, ok := responseBody["pw"]

	assert.False(t, ok)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, loggedIn.ID, player.ID)
	assert.Equal(t, loggedIn.Email, player.Email)
	assert.Equal(t, loggedIn.Username, player.Username)
	assert.Equal(t, player.ID, playerIDFromClaims)
	assert.Equal(t, player.CreatedAt.Unix(), loggedIn.CreatedAt.Unix())
	assert.Equal(t, player.UpdatedAt.Unix(), loggedIn.CreatedAt.Unix())
}
