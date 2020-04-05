package handlers

import (
	"github.com/bnkamalesh/webgo"
	"net/http"
)

// Health is a basic health check handler
func Health(w http.ResponseWriter, r *http.Request) {
	webgo.R200(w, "Success")
}
