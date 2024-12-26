package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/oOSomnus/transflate/internal/TaskManager/repository"
	"golang.org/x/crypto/bcrypt"
)

/*
Authenticate authenticates a user by comparing the provided password with the stored hashed password.

Parameters:
  - username (string): The username of the user to authenticate.
  - password (string): The plaintext password provided by the user.

Returns:
  - (bool): True if the user is authenticated successfully, false otherwise.
  - (error): An error if the user is not found or other issues occur during authentication.
*/

func Authenticate(username, password string) (bool, error) {
	// Fetch pwd
	pwdHashFromDB, err := repository.FindUsrWithUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("invalid username or password")
		}
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Compare
	if err := bcrypt.CompareHashAndPassword([]byte(pwdHashFromDB), []byte(password)); err != nil {
		return false, errors.New("invalid username or password")
	}

	// If match
	return true, nil
}

/*
CreateUser creates a new user with the specified username and password.

Parameters:
  - username (string): The desired username for the new user.
  - password (string): The desired password for the new user.

Returns:
  - (error): An error if the input is invalid, the username already exists, or if any internal process fails.
*/

func CreateUser(username string, password string) error {
	// Validate input
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}

	// Check if user exists
	exists, err := repository.IfUserExists(username)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return errors.New("user with this username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user
	if err := repository.CreateUser(username, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

/*
DecreaseBalance is a wrapper function that calls the primary DecreaseBalance implementation and handles any errors.

Parameters:
  - username (string): The username of the user whose balance is to be reduced.
  - balance (int): The amount to be deducted from the user's balance.

Returns:
  - (error): An error if the underlying DecreaseBalance call fails, wrapped with additional context.
*/
func DecreaseBalance(username string, balance int) error {
	err := repository.DecreaseBalance(username, balance)
	if err != nil {
		return fmt.Errorf("failed to decrease balance: %w", err)
	}
	return nil
}
