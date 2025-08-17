package auth

import (
	// "fmt"
	"golang.org/x/crypto/bcrypt"
)

const cost = 10

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(passwordBytes, cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(hashedPassword, password string) error {
	passwordBytes, hashBytes := []byte(password), []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
	if err != nil {
		return err
	}

	return nil
}