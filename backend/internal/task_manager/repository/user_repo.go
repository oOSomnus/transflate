package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type UserRepository interface {
	FindUsrWithUsername(username string) (string, error)
	IfUserExists(username string) (bool, error)
	CreateUser(username string, password string) error
	DecreaseBalance(username string, balance int) error
	GetBalance(username string) (int, error)
}

type UserRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB: db,
	}
}

// FindUsrWithUsername retrieves the password hash associated with a given username from the database.
// Returns the password hash or an error if the user does not exist or a query error occurs.
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

// IfUserExists checks if a user exists in the database based on the provided username.
// Returns true and nil if the user exists, false and nil if the user does not exist,
// or false and an error if any database error occurs.
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

// CreateUser inserts a new user into the database with the specified username and password.
// Returns an error if the database operation fails.
func (r *UserRepositoryImpl) CreateUser(username string, password string) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err := r.DB.Exec(query, username, password)
	if err != nil {
		return err
	}
	return nil
}

// DecreaseBalance decreases the user's balance by the specified amount, ensuring atomicity and handling possible errors.
// Parameters:
// - username (string): The username of the user whose balance is to be reduced.
// - balance (int): The amount to be deducted from the user's balance.
// Returns:
// - error: An error if the operation fails due to invalid input, insufficient balance, or database operation issues.
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

// GetBalance retrieves the account balance for the specified username from the database.
// Returns the balance as an integer and an error if the query fails or the user is not found.
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
