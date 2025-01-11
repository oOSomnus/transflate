package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
)

/*
AuthMiddleware is a Gin middleware function for handling JWT-based authentication.

Returns:
  - (gin.HandlerFunc): A middleware function to be used in the Gin router.
*/
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const (
			errMissingAuthHeader    = "Authorization header is required"
			errInvalidBearerToken   = "Bearer token required"
			errInvalidToken         = "Invalid or expired token"
			errInvalidTokenClaims   = "Invalid token claims"
			unexpectedSigningMethod = "unexpected signing method"
		)

		authorizationHeader := getAuthorizationHeader(c)
		if authorizationHeader == "" {
			abortWithError(c, http.StatusUnauthorized, errMissingAuthHeader)
			return
		}

		jwtToken := extractBearerToken(authorizationHeader)
		if jwtToken == "" {
			abortWithError(c, http.StatusUnauthorized, errInvalidBearerToken)
			return
		}

		token, err := validateToken(jwtToken)
		if err != nil || !token.Valid {
			log.Printf("Invalid token: %v", err)
			abortWithError(c, http.StatusUnauthorized, errInvalidToken)
			return
		}

		if claims, ok := extractClaims(token); ok {
			c.Set("username", claims["username"])
		} else {
			abortWithError(c, http.StatusUnauthorized, errInvalidTokenClaims)
			return
		}

		log.Println("Authentication successful.")
		c.Next()
	}
}

// getAuthorizationHeader retrieves the "Authorization" header from the provided Gin context and returns its value.
func getAuthorizationHeader(c *gin.Context) string {
	log.Println("Retrieving Authorization header...")
	return c.GetHeader("Authorization")
}

// extractBearerToken extracts the JWT token from the Authorization header by removing the "Bearer " prefix.
func extractBearerToken(authHeader string) string {
	log.Println("Extracting JWT token from header...")
	return strings.TrimPrefix(authHeader, "Bearer ")
}

// validateToken parses and validates a JWT token string for authenticity and signature using the configured secret key.
// Returns the parsed token and an error if the token is invalid or uses an unexpected signing method.
func validateToken(tokenString string) (*jwt.Token, error) {
	log.Printf("Validating token: %s", tokenString)
	secretKey := []byte(viper.GetString("jwt.secret"))
	return jwt.Parse(
		tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secretKey, nil
		},
	)
}

// extractClaims extracts the claims from a given JWT token and checks if they are of type jwt.MapClaims.
// Returns the claims and a boolean indicating success.
func extractClaims(token *jwt.Token) (jwt.MapClaims, bool) {
	log.Println("Extracting token claims...")
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}

// abortWithError sends a JSON response with the specified status code and error message, then aborts further processing.
func abortWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
