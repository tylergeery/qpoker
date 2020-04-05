package app

import (
	"qpoker/http/handlers"

	"github.com/gofiber/fiber"
)

// CreateApp return a new fiber app
func CreateApp() *fiber.App {
	// Start HTTP
	app := fiber.New()

	// Static Asset Routing

	// Web Routing
	// app.Get("/login", handlers.PageLogin)
	// app.Get("/profile", handlers.PageProfile)
	// app.Get("/signup", handlers.PageSignup)
	// app.Get("/game/:id/settings", middleware.AuthorizeAndSetRedirect, handlers.PageGameSettings)
	// app.Get("/:token", middleware.AuthorizeAndSetRedirect, handlers.PageJoin)

	// API Routing
	app.Post("/api/v1/players", handlers.CreatePlayer)
	// app.Post("/api/v1/players/login", handlers.LoginPlayer)
	// app.Get("/api/v1/players/:id", middleware.Authorize, handlers.GetPlayer)
	// app.Post("/api/v1/games", middleware.Authorize, handlers.CreateGame)
	// app.Post("/api/v1/games/:id/join", middleware.Authorize, handlers.CreateGame)
	// app.Get("/api/v1/games/:id", middleware.Authorize, handlers.GetGame)

	return app
}
