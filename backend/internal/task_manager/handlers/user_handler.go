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

// UserHandler defines methods for handling user-related HTTP requests.
// Login handles user authentication requests.
// Register handles user registration requests.
// Info retrieves information about the authenticated user.
type UserHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	Info(c *gin.Context)
}

// UserHandlerImpl handles HTTP requests related to user operations, delegating logic to the associated UserUsecase.
type UserHandlerImpl struct {
	Usecase usecase.UserUsecase
}

// NewUserHandler initializes and returns a new instance of UserHandlerImpl with the provided UserUsecase.
func NewUserHandler(u usecase.UserUsecase) *UserHandlerImpl {
	return &UserHandlerImpl{Usecase: u}
}

// bindJSONAndValidate binds JSON data from the request to a specified struct and validates it.
// Returns false if binding or validation fails, with a corresponding HTTP 400 response.
// Returns true if the process succeeds without errors.
func (h *UserHandlerImpl) bindJSONAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errInvalidRequest})
		return false
	}
	return true
}

// Login handles user login requests by validating credentials and generating a JWT token upon successful authentication.
func (h *UserHandlerImpl) Login(c *gin.Context) {
	var userRequest domain.UserRequest

	if !h.bindJSONAndValidate(c, &userRequest) {
		return
	}

	if err := utils.VerifyTurnstileToken(userRequest.TurnstileToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errTurnstileFailure})
		return
	}

	isAuthenticated, err := h.Usecase.Authenticate(userRequest.Username, userRequest.Password)
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

// Register handles the registration of a new user by validating the input and creating a user through the UserUsecase.
func (h *UserHandlerImpl) Register(c *gin.Context) {
	var userRequest domain.UserRequest

	if !h.bindJSONAndValidate(c, &userRequest) {
		log.Println("Error: Invalid registration request")
		return
	}

	if err := h.Usecase.CreateUser(userRequest.Username, userRequest.Password); err != nil {
		// Hash the user's password before saving it to the database
		hashedPassword, err := utils.HashPassword(userRequest.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			log.Printf("Error: Failed to hash password for username %s: %v", userRequest.Username, err)
			return
		}

		// Use the hashed password when creating the user
		if err := h.Usecase.CreateUser(userRequest.Username, hashedPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			log.Printf("Error: Failed to create user for username %s: %v", userRequest.Username, err)
			return
		}

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

// Info handles authenticated requests to retrieve user information, including username and account balance.
func (h *UserHandlerImpl) Info(c *gin.Context) {
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

	balance, err := h.Usecase.CheckBalance(usernameStr)
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
