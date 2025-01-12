package config

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/spf13/viper"
	"log"
	"net/url"
)

// DatabaseConfig defines an interface for establishing a connection to a database and returning an SQL database handle.
type DatabaseConfig interface {
	Connect() (*sql.DB, error)
}

// DatabaseConfigImpl represents the concrete implementation of configuration settings required to connect to a database.
type DatabaseConfigImpl struct {
}

// PostgresConfig represents the configuration required to establish a connection to a PostgreSQL database.
type PostgresConfig struct {
}

// NewPostgresConfig initializes and returns a new instance of the PostgresConfig struct.
func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{}
}

// Connect establishes a connection to a PostgreSQL database using settings from the configuration and returns the DB instance.
// If the connection or ping fails, the function logs the error and exits the application.
func (p *PostgresConfig) Connect() (*sql.DB, error) {
	// DB connection info
	username := url.QueryEscape(viper.GetString("pg.username"))
	password := url.QueryEscape(viper.GetString("pg.password"))
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
	return db, nil
}
