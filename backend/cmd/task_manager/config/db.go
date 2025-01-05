package config

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/oOSomnus/transflate/pkg/utils"
	"log"
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
	utils.LoadEnv()
	username := utils.GetEnv("PG_USERNAME")
	password := utils.GetEnv("PG_PASSWORD")
	host := utils.GetEnv("PG_HOST")
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
