package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("cannot hashed password %s", err)
	}

	return string(hashPassword), nil
}

func CheckPassword(password string, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
