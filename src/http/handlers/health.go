package handlers

import (
	"github.com/gofiber/fiber"
)

// Health is a basic health check handler
func Health(c *fiber.Ctx) {
	c.JSON(fiber.Map{
		"status": "success",
	})
}
