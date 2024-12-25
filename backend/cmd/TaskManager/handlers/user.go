package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/services/TaskManager/domain"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	isAuthenticated, err := domain.Authenticate(req.Username, req.Password)
	if !isAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "username": req.Username})
}
