package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

/*
LoadEnv loads environment variables from a .env file into the application.

Behavior:
  - Attempts to load variables from a .env file in the current working directory.
  - Logs a fatal error and terminates the application if the .env file cannot be loaded.

Parameters:
  - None

Returns:
  - None
*/
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

/*
GetEnv retrieves the value of an environment variable by its key.

Parameters:
  - key (string): The name of the environment variable to retrieve.

Returns:
  - (string): The value of the environment variable.

Behavior:
  - Logs a fatal error and terminates the application if the environment variable is not set.
*/
func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}
