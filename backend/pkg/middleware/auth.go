package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
AuthMiddleware is a Gin middleware function for handling JWT-based authentication.

Returns:
  - (gin.HandlerFunc): A middleware function to be used in the Gin router.
*/
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization Header
		log.Println("Getting authHeader ...")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		log.Println("Parsing authHeader ...")
		// check prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}
		log.Println("Validating token ...")
		// parse and validate Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// check sign method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			secretKey := []byte(os.Getenv("JWT_SECRET"))
			return secretKey, nil
		})

		// check if valid
		if err != nil || !token.Valid {
			log.Println("Token is invalid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// fetch claim and store in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("username", claims["username"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		log.Println("Auth checked successfully.")
		// proceed
		c.Next()
	}
}
