package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"qpoker/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGameTypes(t *testing.T) {
	// Given
	var retrievedGameTypes []GameType

	player := test.CreateTestPlayer()
	headers := map[string]string{"Content-Type": "application/json", "Authorization": fmt.Sprintf("Bearer %s", player.Token)}
	req := test.CreateTestRequest("GET", "/api/v1/games/types", headers, nil)
	app := CreateApp()

	// When
	response, err := app.Test(req)
	assert.NoError(t, err)

	// Then
	content, err := ioutil.ReadAll(response.Body)
	assert.NoError(t, err)

	json.Unmarshal(content, &retrievedGameTypes)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, 2, len(retrievedGameTypes))
	assert.Equal(t, int64(2), retrievedGameTypes[0].ID)
	assert.Equal(t, int64(1), retrievedGameTypes[1].ID)
	assert.Equal(t, "Hearts", retrievedGameTypes[0].DisplayName)
	assert.Equal(t, "Texas Holdem (Poker)", retrievedGameTypes[1].DisplayName)
	assert.Equal(t, 3, len(retrievedGameTypes[0].Options))
	assert.Equal(t, 6, len(retrievedGameTypes[1].Options))
}
