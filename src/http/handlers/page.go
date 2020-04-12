package handlers

import (
	"qpoker/models"
	"strings"

	"github.com/gofiber/fiber"
)

// PageLanding renders the default app landing page
func PageLanding(c *fiber.Ctx) {
	c.Render("main.mustache", fiber.Map{})
}

// PageTable renders a poker table
func PageTable(c *fiber.Ctx) {
	gameSlug := strings.ToLower(c.Params("slug"))

	_, err := models.GetGameBy("slug", gameSlug)
	if err != nil {
		c.SendStatus(404)
		c.Render("error.mustache", fiber.Map{})
		return
	}

	c.Render("table.mustache", fiber.Map{})
}
