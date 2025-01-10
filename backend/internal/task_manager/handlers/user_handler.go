package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/internal/task_manager/domain"
	"github.com/oOSomnus/transflate/internal/task_manager/usecase"
	"github.com/oOSomnus/transflate/pkg/utils"
)

const (
	errInvalidRequest     = "Invalid request"
	errTurnstileFailure   = "Turnstile token verification failed"
	errAuthFailure        = "Failed to authenticate, please check the user information"
	errTokenGeneration    = "Failed to generate token"
	errUserUnauthorized   = "User unauthorized"
	errInvalidUsername    = "Invalid username"
	errBalanceCheckFailed = "Error checking balance"
)

// bindJSONAndValidate is a helper function to bind and validate JSON request payloads.
func bindJSONAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidRequest})
		return false
	}
	return true
}

// Login handles user authentication by validating credentials, verifying Turnstile token, and returning a JWT token.
func Login(c *gin.Context) {
	var userRequest domain.UserRequest

	if !bindJSONAndValidate(c, &userRequest) {
		return
	}

	if err := utils.VerifyTurnstileToken(userRequest.TurnstileToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errTurnstileFailure})
		return
	}

	isAuthenticated, err := usecase.Authenticate(userRequest.Username, userRequest.Password)
	if err != nil || !isAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errAuthFailure})
		return
	}

	token, err := utils.GenerateToken(userRequest.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errTokenGeneration})
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message":  "Login successful",
			"username": userRequest.Username,
			"token":    token,
		},
	)
}

// Register handles user registration by processing the incoming JSON request, validating input, and creating a new user.
func Register(c *gin.Context) {
	var userRequest domain.UserRequest

	if !bindJSONAndValidate(c, &userRequest) {
		log.Println("Error: Invalid registration request")
		return
	}

	if err := usecase.CreateUser(userRequest.Username, userRequest.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		log.Printf("Error: Failed to create user for username %s: %v", userRequest.Username, err)
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"message": "User created successfully",
		},
	)
}

// Info handles retrieving user information, including username and balance, and returns JSON responses based on the outcome.
func Info(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errUserUnauthorized})
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidUsername})
		return
	}

	balance, err := usecase.CheckBalance(usernameStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errBalanceCheckFailed})
		log.Printf("Error: Failed to check balance for username %s: %v", usernameStr, err)
		return
	}

	c.JSON(
		http.StatusOK, gin.H{
			"username": usernameStr,
			"balance":  balance,
		},
	)
}
