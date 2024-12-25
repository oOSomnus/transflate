package repository

import (
	"database/sql"
	"errors"
	"github.com/oOSomnus/transflate/cmd/TaskManager/config"
)

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
