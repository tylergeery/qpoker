package handlers

import (
	"github.com/gofiber/fiber"
)

// PageLanding renders the default app landing page
func PageLanding(ctx *fiber.Ctx) {
	ctx.Render("main.mustache", fiber.Map{})
}

// PageTable renders a poker table
func PageTable(ctx *fiber.Ctx) {
	playerID := c.Locals("playerID")

	ctx.Render("table.mustache", fiber.Map{})
}
