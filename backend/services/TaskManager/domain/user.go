package domain

import (
	"database/sql"
	"errors"

	"github.com/oOSomnus/transflate/cmd/TaskManager/config"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func Authenticate(username, password string) (bool, error) {
	var user User

	query := "SELECT id, username, password FROM users WHERE username = ? AND password = ?"
	row := config.DB.QueryRow(query, username, password)

	// 如果找到匹配用户
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("invalid username or password")
		}
		return false, err
	}

	return true, nil
}
