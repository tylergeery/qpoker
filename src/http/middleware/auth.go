package middleware

import (
	"fmt"
	"qpoker/auth"
	"strings"
	"time"

	"github.com/gofiber/fiber"
)

// Authorize ensures valid player tokens
func Authorize(c *fiber.Ctx) {
	authHeader := c.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	playerID, err := auth.GetPlayerIDFromAccessToken(token)
	if err != nil {
		fmt.Printf("Player token parse error: %s\n", err)
		c.SendString(fmt.Sprintf("Forbidden: %s", err))
		c.SendStatus(403)
	}

	// add playerID to context
	c.Locals("playerID", playerID)
	c.Next(nil)
}

// AuthorizeAndSetRedirect sets a redirect before Authorize
func AuthorizeAndSetRedirect(c *fiber.Ctx) {
	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "authorize_redirect"
	cookie.Value = c.Path()
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.Cookie(cookie)

	Authorize(c)
}
