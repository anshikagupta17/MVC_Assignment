package models

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPass(pass string) (string, error) {
	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed_pass), nil
}
