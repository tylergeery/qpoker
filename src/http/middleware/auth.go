package middleware

import (
	"fmt"
	"qpoker/auth"
	"qpoker/http/utils"
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
		fmt.Printf("Player token (%s) parse error: %s\n", token, err)
		c.SendStatus(403)
		c.JSON(utils.FormatErrors(err))
		return
	}

	// add playerID to context
	c.Locals("playerID", playerID)
	c.Locals("token", token)
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
