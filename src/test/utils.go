package test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"qpoker/auth"
	"qpoker/models"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// TestPass hold the test password for tests in need of knowing it
var TestPass = "testpw"

// CreateTestPlayer creates a new player for a test
func CreateTestPlayer() *models.Player {
	ts := time.Now().UTC().UnixNano()
	player := &models.Player{
		Username: fmt.Sprintf("testplayer_%d", ts),
		Email:    fmt.Sprintf("testplayer_%d@test.com", ts),
	}
	_ = player.Create(TestPass)
	player.Token, _ = auth.CreateToken(player)

	return player
}

// CreateTestGame creates a new game for a test
func CreateTestGame(player *models.Player, gameTypeID int64) *models.Game {
	ts := time.Now().UTC().UnixNano()
	game := &models.Game{
		Name:       fmt.Sprintf("Test Game %d", ts),
		OwnerID:    player.ID,
		Status:     models.GameStatusInit,
		GameTypeID: gameTypeID,
	}
	_ = game.Save()
	options, _ := models.GetGameOptionsForGame(game.ID, game.GameTypeID)
	game.Options = options.Options

	return game
}

// CreateTestRequest creates a new test request
func CreateTestRequest(action, endpoint string, headers map[string]string, body map[string]interface{}) *http.Request {
	var request *http.Request

	if body == nil {
		request = httptest.NewRequest(action, fmt.Sprintf("http://qpoker.com%s", endpoint), nil)
	} else {
		content, _ := json.Marshal(body)
		reader := strings.NewReader(string(content))

		request = httptest.NewRequest(action, fmt.Sprintf("http://qpoker.com%s", endpoint), reader)
		request.Header.Set("Content-Length", strconv.Itoa(len(string(content))))
		request.Header.Set("Content-Type", "application/json")
	}

	for key, val := range headers {
		request.Header.Set(key, val)
	}

	return request
}
