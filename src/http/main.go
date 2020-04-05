package main

import (
	"qpoker/http/app"
	"qpoker/models"
)

func main() {
	// Ensure connection to DB
	db := models.ConnectToDB()
	defer db.Close()

	server := app.CreateApp()
	server.Listen(8080)
}
