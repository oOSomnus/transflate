package config

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/spf13/viper"
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
	//utils.LoadEnv()
	username := viper.GetString("pg.username")
	password := viper.GetString("pg.password")
	host := viper.GetString("pg.host")
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
