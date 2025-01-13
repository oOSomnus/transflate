package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// UserRepository defines methods for user management including retrieval, creation, and balance operations.
// FindUsrWithUsername retrieves a user's information using their username.
// IfUserExists checks whether a user with the given username exists.
// CreateUser creates a new user with the specified username and password.
// DecreaseBalance reduces the user's balance by the specified amount.
// GetBalance retrieves the current balance of the user.
type UserRepository interface {
	FindUsrWithUsername(username string) (string, error)
	IfUserExists(username string) (bool, error)
	CreateUser(username string, password string) error
	DecreaseBalance(username string, balance int) error
	GetBalance(username string) (int, error)
}

// UserRepositoryImpl interacts with the database to perform user-related operations like querying and updating data.
// It wraps around an *sql.DB instance for executing SQL queries and managing transactions.
type UserRepositoryImpl struct {
	DB *sql.DB
}

// NewUserRepository initializes a new UserRepositoryImpl with a given sql.DB connection and returns its instance.
func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB: db,
	}
}

// FindUsrWithUsername retrieves the hashed password for a given username from the database or returns an error if not found.
func (r *UserRepositoryImpl) FindUsrWithUsername(username string) (string, error) {
	query := "SELECT password FROM users WHERE username = $1"
	row := r.DB.QueryRow(query, username)
	var pwd string
	err := row.Scan(&pwd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return pwd, nil
}

// IfUserExists checks if a user with the given username exists in the database. Returns true if user exists, else false.
func (r *UserRepositoryImpl) IfUserExists(username string) (bool, error) {
	query := "SELECT userid FROM users WHERE username = $1"
	row := r.DB.QueryRow(query, username)
	var userId int
	err := row.Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateUser inserts a new user into the database with the provided username and password.
// It returns an error if the operation fails.
func (r *UserRepositoryImpl) CreateUser(username string, password string) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err := r.DB.Exec(query, username, password)
	if err != nil {
		return err
	}
	return nil
}

// DecreaseBalance decreases the balance of the specified user by the given amount.
// Returns an error if the balance is insufficient, the user is not found, or there is a transaction/database failure.
func (r *UserRepositoryImpl) DecreaseBalance(username string, balance int) error {
	if balance <= 0 {
		return errors.New("invalid amount")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			log.Println("recovered from panic:", p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	var currentBalance int
	query := "SELECT balance FROM users WHERE username = $1 FOR UPDATE"
	err = tx.QueryRowContext(ctx, query, username).Scan(&currentBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get current balance: %w", err)
	}

	if currentBalance < balance {
		return errors.New("insufficient balance")
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = balance - $1 WHERE username = $2", balance, username)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}

// GetBalance retrieves the balance of a user by their username from the database.
// It returns the balance as an integer or an error if an issue occurs or the user is not found.
func (r *UserRepositoryImpl) GetBalance(username string) (int, error) {
	query := "SELECT balance FROM users WHERE username = $1"
	row := r.DB.QueryRow(query, username)
	var balance int
	err := row.Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("user not found")
		}
		return 0, err
	}
	return balance, nil
}
