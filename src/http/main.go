package main

import (
	"qpoker/http/handlers"
	"qpoker/models"
)

func main() {
	// Ensure connection to DB
	db := models.ConnectToDB()
	defer db.Close()

	app := handlers.CreateApp()
	app.Listen(8080)
}
