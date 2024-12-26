package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/internal/TaskManager/domain"
	"github.com/oOSomnus/transflate/internal/TaskManager/usecase"
	"github.com/oOSomnus/transflate/pkg/utils"
	"net/http"
)

/*
Login handles user login by authenticating credentials and generating a JWT token.

Parameters:
  - c (*gin.Context): The Gin context containing the HTTP request and response objects.

Returns:
  - Responds directly to the HTTP client with appropriate status codes and messages based on the operation's success or failure.
*/

func Login(c *gin.Context) {
	var req domain.UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	isAuthenticated, err := usecase.Authenticate(req.Username, req.Password)
	if !isAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT Token
	token, err := utils.GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"username": req.Username,
		"token":    token,
	})
}

/*
Register handles the HTTP request to register a new user.

Parameters:
  - c (*gin.Context): The context of the current HTTP request, providing request and response handling.

Returns:
  - (JSON): A success or error message depending on the outcome.
*/
func Register(c *gin.Context) {
	var req domain.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := usecase.CreateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}
