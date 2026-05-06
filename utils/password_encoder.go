package utils

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {

	hashedpassword, bcryptErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if bcryptErr != nil {
		log.Printf("ERRO: Falha ao encriptar senha! Message: %s", bcryptErr.Error())
		return "", errors.New("Erro Interno!")
	}

	return string(hashedpassword), nil
}
