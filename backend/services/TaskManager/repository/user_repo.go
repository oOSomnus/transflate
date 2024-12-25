package repository

import (
	"database/sql"
	"errors"
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
	query := "SELECT password FROM users WHERE username = ?"
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
	query := "SELECT userid FROM users WHERE username = ?"
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
	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err := config.DB.Exec(query, username, password)
	if err != nil {
		return err
	}
	return nil
}
