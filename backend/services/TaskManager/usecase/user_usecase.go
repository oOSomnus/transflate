package usecase

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/oOSomnus/transflate/services/TaskManager/repository"
)

func Authenticate(username, password string) (bool, error) {

	// Fetch hashed val
	pwdHashFromDB, err := repository.FindUsrWithUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("user not found")
		}
		return false, err
	}

	// Compute hash of input
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedPwdInput := hex.EncodeToString(hash.Sum(nil))

	// Compare hashed value
	if hashedPwdInput == pwdHashFromDB {
		return true, nil
	}

	return false, errors.New("invalid password")
}
