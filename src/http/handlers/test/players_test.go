package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"qpoker/auth"
	"qpoker/http/app"
	"qpoker/http/test"
	"qpoker/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPlayerCreateInvalid(t *testing.T) {
	type TestCase struct {
		body map[string]string
		expected string
	}
	cases := []TestCase{
		TestCase{
			body: map[string]string{
				"email": "test@qpoker.com",
				"pw": "testpassword",
			},
			expected: "{\"errors\": [\"username requires at least 6 characters\"]}",
		},
		TestCase{
			body: map[string]string{
				"username": "testqpokerguy",
				"pw": "testpassword",
			},
			expected: "{\"errors\": [\"Invalid email format: \"]}",
		},
		TestCase{
			body: map[string]string{
				"email": "akslfadfalasdfas",
				"username": "testqpokerguy",
				"pw": "testpassword",
			},
			expected: "{\"errors\": [\"Invalid email format: akslfadfalasdfas\"]}",
		},
		TestCase{
			body: map[string]string{
				"email": "test@qpoker.com",
				"username": "testqpokerguy",
			},
			expected: "{\"errors\": [\"password requires at least 6 characters\"]}",
		},
	}

	server := app.CreateApp()

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
		"email": email,
		"username": username,
		"pw": "testpasstest",
	}
	req := test.CreateTestRequest("POST", "/api/v1/players", nil, body)
	server := app.CreateApp()

	// When
	response, err := server.Test(req)
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
