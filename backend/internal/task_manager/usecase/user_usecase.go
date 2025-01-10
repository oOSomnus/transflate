package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/oOSomnus/transflate/internal/task_manager/repository"
	"golang.org/x/crypto/bcrypt"
)

// 错误消息常量
const (
	ErrInvalidCredentials = "invalid username or password"
	ErrUserAlreadyExists  = "user already exists"
	ErrEmptyInput         = "username and password cannot be empty"
)

// Authenticate validates credentials by comparing a provided password to the stored hash.
// Returns true if valid, otherwise returns false and an error.
func Authenticate(username, password string) (bool, error) {
	hashedPassword, err := repository.FindUsrWithUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New(ErrInvalidCredentials)
		}
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if err := validatePassword(password, hashedPassword); err != nil {
		return false, err
	}
	return true, nil
}

// validatePassword compares the plaintext password with the hashed password.
// Returns an error if passwords don't match.
func validatePassword(password, hashedPassword string) error {
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return errors.New(ErrInvalidCredentials)
	}
	return nil
}

// CreateUser creates a new user with a hashed password.
// Returns an error if the input is invalid, user already exists, or creation fails.
func CreateUser(username, password string) error {
	if err := validateUserInput(username, password); err != nil {
		return err
	}
	if exists, err := repository.IfUserExists(username); err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	} else if exists {
		return errors.New(ErrUserAlreadyExists)
	}
	return createUserInRepository(username, password)
}

// validateUserInput ensures username and password are not empty.
func validateUserInput(username, password string) error {
	if username == "" || password == "" {
		return errors.New(ErrEmptyInput)
	}
	return nil
}

// createUserInRepository hashes the password and adds the user to the repository.
func createUserInRepository(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	if err = repository.CreateUser(username, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// DecreaseBalance decreases the given user's balance.
// Returns an error if the operation fails.
func DecreaseBalance(username string, balance int) error {
	if err := repository.DecreaseBalance(username, balance); err != nil {
		return fmt.Errorf("failed to decrease balance: %w", err)
	}
	return nil
}

// CheckBalance retrieves and validates a user's balance.
// Returns the balance and an error if the retrieval fails.
func CheckBalance(username string) (int, error) {
	balance, err := repository.GetBalance(username)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}
	if balance < 0 {
		return 0, errors.New("negative balance")
	}
	return balance, nil
}
