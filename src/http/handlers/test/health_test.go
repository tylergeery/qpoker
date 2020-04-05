package test

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"qpoker/http/app"
	"qpoker/http/test"
	"testing"
)

func TestHealth(t *testing.T) {
	server := app.CreateApp()
	req := test.CreateTestRequest("GET", "/health", nil, nil)

	response, err := server.Test(req)
	body, _ := ioutil.ReadAll(response.Body)

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.JSONEq(t, "{\"status\": \"success\"}", string(body))
}
