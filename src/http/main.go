package main

import (
	"math/rand"
	"qpoker/http/handlers"
	"qpoker/models"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Ensure connection to DB
	db := models.ConnectToDB()
	defer db.Close()

	app := handlers.CreateApp()

	app.Listen(8080)
}
