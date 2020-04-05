package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"qpoker/models"
	"strings"
	"time"
)

// CreateTestPlayer creates a new player for a test
func CreateTestPlayer() *models.Player {
	ts := time.Now().UTC().UnixNano()
	player := &models.Player{
		Username: fmt.Sprintf("testplayer_%s", ts)
		Email: fmt.Sprintf("testplayer_%s@test.com", ts),
	}
	_ = player.Create("testpw")

	return player
}

// CreateTestRequest creates a new test request
func CreateTestRequest(action, endpoint string, headers, body map[string]string) *http.Request {
	var request *http.Request
	if body == nil {
		request = httptest.NewRequest(action, fmt.Sprintf("http://qpoker.com/api/v1/%s", endpoint), nil)
	} else {
		content, _ := json.Marshal(body)
		reader := strings.NewReader(string(content))
		request = httptest.NewRequest(action, fmt.Sprintf("http://qpoker.com/api/v1/%s", endpoint), reader)
	}

	for key, val := range headers {
		request.Header.Set(key, val)
	}

	return request
}
