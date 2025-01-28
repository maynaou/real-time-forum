package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		fmt.Println("Failed to hash password:", err)
		return nil, err
	}

	return HashedPassword, nil
}

func ComparePasswords(hashedPassword, plaintextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))

	if err != nil {
		fmt.Println("Password comparison failed")
		return false
	}

	return true
}
