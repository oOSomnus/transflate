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

// AuthMiddleware is a middleware function that verifies the validity of a JWT token in the Authorization header.
// It ensures the token is present, well-formed, and contains valid claims.
// If the token is invalid, expired, or missing, the request is aborted with an appropriate HTTP status and error message.
// Upon successful validation, the claims (e.g., "username") are extracted and added to the Gin context for downstream handlers.
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

// getAuthorizationHeader retrieves the Authorization header value from the provided gin.Context.
func getAuthorizationHeader(c *gin.Context) string {
	log.Println("Retrieving Authorization header...")
	return c.GetHeader("Authorization")
}

// extractBearerToken extracts the Bearer token from the given Authorization header string.
func extractBearerToken(authHeader string) string {
	log.Println("Extracting JWT token from header...")
	return strings.TrimPrefix(authHeader, "Bearer ")
}

// validateToken parses and validates a JWT token string using a predefined secret key and signing method.
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

// extractClaims extracts claims from a JWT token and asserts them as jwt.MapClaims.
// Returns the extracted claims and a boolean indicating the success of the assertion.
func extractClaims(token *jwt.Token) (jwt.MapClaims, bool) {
	log.Println("Extracting token claims...")
	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}

// abortWithError sends a JSON response containing an error message with the specified status code and aborts the request.
func abortWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
