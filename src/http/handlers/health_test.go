package handlers

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"qpoker/test"
	"testing"
)

func TestHealth(t *testing.T) {
	server := CreateApp()
	req := test.CreateTestRequest("GET", "/health", nil, nil)

	response, err := server.Test(req)
	body, _ := ioutil.ReadAll(response.Body)

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.JSONEq(t, "{\"status\": \"success\"}", string(body))
}
