package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/oOSomnus/transflate/cmd/TaskManager/config"
)

/*
FindUsrWithUsername retrieves the hashed password for a user with the given username.

Parameters:
  - username (string): The username of the user whose password is being retrieved.

Returns:
  - (string): The hashed password associated with the given username.
  - (error): An error if the user is not found or any database issues occur.
*/
func FindUsrWithUsername(username string) (string, error) {
	query := "SELECT password FROM users WHERE username = $1"
	row := config.DB.QueryRow(query, username)
	var pwd string
	err := row.Scan(&pwd)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}
	return pwd, nil
}

/*
IfUserExists checks if a user with the given username exists in the database.

Parameters:
  - username (string): The username of the user to check.

Returns:
  - (bool): True if the user exists, false otherwise.
  - (error): An error if there are database-related issues.
*/
func IfUserExists(username string) (bool, error) {
	query := "SELECT userid FROM users WHERE username = $1"
	row := config.DB.QueryRow(query, username)
	var userId int
	err := row.Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

/*
CreateUser creates a new user by inserting their username and hashed password into the database.

Parameters:
  - username (string): The username of the new user to be created.
  - password (string): The hashed password of the user to be stored.

Returns:
  - (error): Returns an error if the insertion into the database fails, or nil if the operation is successful.
*/
func CreateUser(username string, password string) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err := config.DB.Exec(query, username, password)
	if err != nil {
		return err
	}
	return nil
}

/*
DecreaseBalance reduces the balance of a user in the database by a specified amount.

Parameters:
  - username (string): The username of the user whose balance is to be reduced.
  - balance (int): The amount to be deducted from the user's balance. Must be greater than zero.

Returns:
  - (error): An error if any of the following occur:
  - The balance parameter is invalid (less than or equal to zero).
  - A transaction cannot be initiated.
  - The user is not found in the database.
  - The user has insufficient balance.
  - There is an error executing the database query to update the balance.
*/
func DecreaseBalance(username string, balance int) error {
	if balance <= 0 {
		return errors.New("invalid amount")
	}
	tx, err := config.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		// rollback
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()
	var currentBalance int
	query := "SELECT balance FROM users WHERE username = $1 FOR UPDATE"
	err = tx.QueryRowContext(context.Background(), query, username).Scan(&currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get current balance: %w", err)
	}
	if currentBalance < balance {
		return errors.New("insufficient balance")
	}
	_, err = tx.ExecContext(context.Background(), "UPDATE users SET balance = balance - $1 WHERE username = $2", balance, username)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}
