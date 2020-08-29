package handlers

import (
	"qpoker/http/middleware"

	"github.com/gofiber/fiber"
)

// CreateApp return a new fiber app
func CreateApp() *fiber.App {
	// Start HTTP
	app := fiber.New()

	// Health
	app.Get("/health", Health)

	// Static Asset Routing

	// Web Routing
	app.Get("/", PageLanding)
	app.Get("/:slug", PageTable)

	// API Routing
	apiV1 := app.Group("/api/v1")
	apiV1.Post("/players", CreatePlayer)
	apiV1.Post("/players/login", LoginPlayer)
	apiV1.Put("/players/:id", middleware.Authorize, UpdatePlayer)
	apiV1.Get("/players/:id", middleware.Authorize, GetPlayer)

	apiV1.Get("/games/types", middleware.Authorize, GetGameTypes)
	apiV1.Post("/games", middleware.Authorize, CreateGame)
	apiV1.Post("/games/:slug/join", middleware.Authorize, CreateGame)
	apiV1.Put("/games/:gameID", middleware.Authorize, UpdateGame)
	apiV1.Get("/games/:gameID", middleware.Authorize, GetGame)

	apiV1.Get("/games/:gameID/history", middleware.Authorize, GetGameHistory)

	return app
}
