package models

import (
	"database/sql"
	"log"
	"os"
	"sync"

	// wrapper around sql package for postgres
	_ "github.com/lib/pq"
)

var db *sql.DB
var persistentOnce sync.Once

// ConnectToDB returns a database connection
func ConnectToDB() *sql.DB {
	persistentOnce.Do(func() {
		db = connectToDB()
	})

	return db
}

func connectToDB() *sql.DB {
	connURL := os.Getenv("PG_CONNECTION")
	if connURL == "" {
		log.Fatalf("Failed to find postgress connection value")
	}

	db, err := sql.Open("postgres", connURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB via %s: %v", connURL, err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB via %s: %v", connURL, err)
	}

	log.Println("Connected to DB")

	return db
}
