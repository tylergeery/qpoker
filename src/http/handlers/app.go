package handlers

import (
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
	// app.Get("/login", handlers.PageLogin)
	// app.Get("/profile", handlers.PageProfile)
	// app.Get("/signup", handlers.PageSignup)
	// app.Get("/game/:id/settings", middleware.AuthorizeAndSetRedirect, handlers.PageGameSettings)
	// app.Get("/:token", middleware.AuthorizeAndSetRedirect, handlers.PageJoin)

	// API Routing
	apiV1 := app.Group("/api/v1")
	apiV1.Post("/players", CreatePlayer)
	apiV1.Post("/players/login", LoginPlayer)
	// apiV1.Get("/api/v1/players/:id", middleware.Authorize, handlers.GetPlayer)
	// apiV1.Post("/api/v1/games", middleware.Authorize, handlers.CreateGame)
	// apiV1.Post("/api/v1/games/:id/join", middleware.Authorize, handlers.CreateGame)
	// apiV1.Get("/api/v1/games/:id", middleware.Authorize, handlers.GetGame)

	return app
}
