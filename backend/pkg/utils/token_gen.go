package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

/*
GenerateToken generates a JWT token for the given username.

Parameters:
  - username (string): The username to include in the token claims.

Returns:
  - (string): The generated JWT token as a string.
  - (error): An error if the token signing process fails.
*/
func GenerateToken(username string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	// Token created
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expired after 24h
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
