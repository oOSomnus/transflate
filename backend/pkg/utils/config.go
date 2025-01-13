package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// LoadEnv reads the environment variables from a .env file and loads them into the application's runtime environment.
// It terminates the program with a fatal log if the .env file cannot be loaded.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetEnv retrieves the value of the specified environment variable or logs a fatal error if it is not set.
func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}
