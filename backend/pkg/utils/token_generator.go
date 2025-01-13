package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

// GenerateToken generates a JWT token for the provided username with a validity of 24 hours.
// Returns the signed token as a string and an error if token generation fails.
func GenerateToken(username string) (string, error) {
	//LoadEnv()
	var jwtSecret = []byte(viper.GetString("jwt.secret"))
	// Token created
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expired after 24h
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
