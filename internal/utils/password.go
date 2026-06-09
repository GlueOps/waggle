package utils

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func HashPassword(plain string) (string, error) {
	if plain == "" {
		return "", errors.New("password is empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash: %w", err)
	}
	return string(h), nil
}

func VerifyPassword(plain, hash string) bool {
	if plain == "" || hash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
