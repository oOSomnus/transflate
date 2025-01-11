package config

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/spf13/viper"
	"log"
	"net/url"
)

//var DB *sql.DB

type DatabaseConfig interface {
	Connect() (*sql.DB, error)
}

type DatabaseConfigImpl struct {
}

type PostgresConfig struct {
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{}
}

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

//func ConnectDB() {
//	// DB connection info
//	//utils.LoadEnv()
//	username := viper.GetString("pg.username")
//	password := viper.GetString("pg.password")
//	host := viper.GetString("pg.host")
//	port := "5432"
//	dbname := "postgres"
//	sslmode := "disable"
//
//	// DSN
//	dsn := "postgres://" + username + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
//
//	// Open db conn
//	db, err := sql.Open("postgres", dsn)
//	if err != nil {
//		log.Fatalf("Failed to connect to database: %v", err)
//	}
//
//	// test conn
//	if err := db.Ping(); err != nil {
//		log.Fatalf("Database ping failed: %v", err)
//	}
//
//	log.Println("Connected to the PostgreSQL database!")
//	DB = db
//}
