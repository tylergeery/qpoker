package handlers

import (
	"github.com/gofiber/fiber"
)

// PageLanding renders the default app landing page
func PageLanding(c *fiber.Ctx) {
	c.Render("main.mustache", fiber.Map{})
}

// PageTable renders a poker table
func PageTable(c *fiber.Ctx) {
	c.Render("table.mustache", fiber.Map{})
}
