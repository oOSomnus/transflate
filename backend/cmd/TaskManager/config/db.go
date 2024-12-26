package config

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

/*
ConnectDB establishes a connection to a PostgreSQL database using environment variables for authentication.

Returns:
  - Updates the global variable `DB` with the established database connection.
  - Logs errors and exits the application in case of failures.
*/

func ConnectDB() {
	// DB connection info
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	host := "127.0.0.1"
	port := "5432"
	dbname := "postgres"
	sslmode := "disable"

	// DSN
	dsn := "postgres://" + username + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode

	// Open db conn
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// test conn
	if err := db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Connected to the PostgreSQL database!")
	DB = db
}
