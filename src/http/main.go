package main

import (
	"math/rand"
	"qpoker/http/handlers"
	"qpoker/models"
	"time"

	"github.com/gofiber/template"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Ensure connection to DB
	db := models.ConnectToDB()
	defer db.Close()

	app := handlers.CreateApp()

	// Settings
	app.Settings.TemplateEngine = template.Mustache()
	app.Settings.TemplateFolder = "/src/http/views/"

	app.Listen(8080)
}
