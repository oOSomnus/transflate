package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/oOSomnus/transflate/internal/task_manager/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Authenticate(username, password string) (bool, error)
	CreateUser(username, password string) error
	DecreaseBalance(username string, balance int) error
	CheckBalance(username string) (int, error)
}

type UserUsecaseImpl struct {
	Repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		Repo: r,
	}
}

// Error msg const
const (
	ErrInvalidCredentials = "invalid username or password"
	ErrUserAlreadyExists  = "user already exists"
	ErrEmptyInput         = "username and password cannot be empty"
)

// Authenticate validates credentials by comparing a provided password to the stored hash.
// Returns true if valid, otherwise returns false and an error.
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

// validatePassword compares the plaintext password with the hashed password.
// Returns an error if passwords don't match.
func (u *UserUsecaseImpl) validatePassword(password, hashedPassword string) error {
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return errors.New(ErrInvalidCredentials)
	}
	return nil
}

// CreateUser creates a new user with a hashed password.
// Returns an error if the input is invalid, user already exists, or creation fails.
func (u *UserUsecaseImpl) CreateUser(username, password string) error {
	if err := u.validateUserInput(username, password); err != nil {
		return err
	}
	if exists, err := u.Repo.IfUserExists(username); err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	} else if exists {
		return errors.New(ErrUserAlreadyExists)
	}
	return u.createUserInRepository(username, password)
}

// validateUserInput ensures username and password are not empty.
func (u *UserUsecaseImpl) validateUserInput(username, password string) error {
	if username == "" || password == "" {
		return errors.New(ErrEmptyInput)
	}
	return nil
}

// createUserInRepository hashes the password and adds the user to the repository.
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

// DecreaseBalance decreases the given user's balance.
// Returns an error if the operation fails.
func (u *UserUsecaseImpl) DecreaseBalance(username string, balance int) error {
	if err := u.Repo.DecreaseBalance(username, balance); err != nil {
		return fmt.Errorf("failed to decrease balance: %w", err)
	}
	return nil
}

// CheckBalance retrieves and validates a user's balance.
// Returns the balance and an error if the retrieval fails.
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
