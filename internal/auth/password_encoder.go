package auth

import (
	"unicode/utf8"

	"github.com/CriciumaDevJobs/backend/handlers"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {

	//bcrypt encoder trunkates string with more then 72 characters
	if utf8.RuneCountInString(password) > 72 {
		return "", handlers.ErrPasswordToLong
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedpassword), nil
}
