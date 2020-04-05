package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"qpoker/models"
	"strconv"
	"strings"
	"time"
)

// CreateTestPlayer creates a new player for a test
func CreateTestPlayer() *models.Player {
	ts := time.Now().UTC().UnixNano()
	player := &models.Player{
		Username: fmt.Sprintf("testplayer_%d", ts),
		Email:    fmt.Sprintf("testplayer_%d@test.com", ts),
	}
	_ = player.Create("testpw")

	return player
}

// CreateTestRequest creates a new test request
func CreateTestRequest(action, endpoint string, headers, body map[string]string) *http.Request {
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
