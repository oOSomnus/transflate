package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/oOSomnus/transflate/internal/task_manager/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// UserUsecase defines the operations related to user management in the system.
// Authenticate validates a user's credentials and returns authentication success or failure.
// CreateUser registers a new user with the given username and password.
// DecreaseBalance reduces the balance of a specified user by the provided amount.
// CheckBalance retrieves the current balance of a specified user.
type UserUsecase interface {
	Authenticate(username, password string) (bool, error)
	CreateUser(username, password string) error
	DecreaseBalance(username string, balance int) error
	CheckBalance(username string) (int, error)
}

// UserUsecaseImpl is a struct implementing business use cases for users using a UserRepository.
type UserUsecaseImpl struct {
	Repo repository.UserRepository
}

// NewUserUsecase initializes and returns a new instance of UserUsecaseImpl using the provided UserRepository.
func NewUserUsecase(r repository.UserRepository) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		Repo: r,
	}
}

// ErrInvalidCredentials represents an error message for invalid login credentials.
// ErrUserAlreadyExists represents an error message indicating the user already exists.
// ErrEmptyInput represents an error message for empty username or password inputs.
const (
	ErrInvalidCredentials = "invalid username or password"
	ErrUserAlreadyExists  = "user already exists"
	ErrCreateUser         = "failed to create user"
	ErrEmptyInput         = "username and password cannot be empty"
	ErrInvalidRegInfo     = "invalid registration information"
	ErrInvalidUsrInput    = "invalid user input"
)

// Authenticate verifies if the provided username and password are valid and returns a boolean and potential error.
// It retrieves the hashed password from the repository and validates it against the provided password.
func (u *UserUsecaseImpl) Authenticate(username, password string) (bool, error) {
	hashedPassword, err := u.Repo.FindUsrWithUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New(ErrInvalidCredentials)
		}
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if err := u.validatePassword(password, hashedPassword); err != nil {
		return false, err
	}
	return true, nil
}

// validatePassword compares a plaintext password with a hashed password and returns an error if they do not match.
func (u *UserUsecaseImpl) validatePassword(password, hashedPassword string) error {
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return errors.New(ErrInvalidCredentials)
	}
	return nil
}

// CreateUser creates a new user with the provided username and password, ensuring valid input and uniqueness constraints.
func (u *UserUsecaseImpl) CreateUser(username, password string) error {
	if err := u.validateUserInput(username, password); err != nil {
		log.Println("error validating user input")
		return err
	}
	if len(username) < 4 || len(username) > 15 {
		return errors.New(ErrInvalidRegInfo)
	}
	if len(password) < 11 || len(password) > 18 {
		return errors.New(ErrInvalidRegInfo)
	}

	if exists, err := u.Repo.IfUserExists(username); err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	} else if exists {
		log.Println("user already exists")
		return errors.New(ErrUserAlreadyExists)
	}
	if err := u.createUserInRepository(username, password); err != nil {
		log.Println("error creating user")
		return errors.New(ErrCreateUser)
	}
	return nil
}

// validateUserInput validates the username and password input, ensuring they are not empty. Returns an error if invalid.
func (u *UserUsecaseImpl) validateUserInput(username, password string) error {
	if username == "" || password == "" {
		return errors.New(ErrEmptyInput)
	}
	return nil
}

// createUserInRepository hashes the password and saves the user with the hashed password in the repository.
func (u *UserUsecaseImpl) createUserInRepository(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	if err = u.Repo.CreateUser(username, string(hashedPassword)); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// DecreaseBalance decreases the specified balance for a user identified by username and returns an error if the operation fails.
func (u *UserUsecaseImpl) DecreaseBalance(username string, balance int) error {
	if err := u.Repo.DecreaseBalance(username, balance); err != nil {
		return fmt.Errorf("failed to decrease balance: %w", err)
	}
	return nil
}

// CheckBalance retrieves the balance of a specific user by username and returns it.
// It returns an error if the balance retrieval fails or is negative.
func (u *UserUsecaseImpl) CheckBalance(username string) (int, error) {
	balance, err := u.Repo.GetBalance(username)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}
	if balance < 0 {
		return 0, errors.New("negative balance")
	}
	return balance, nil
}
