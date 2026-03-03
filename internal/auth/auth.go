package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)

// hash the password using
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Printf("Error obtaining hash: %v", err)
		return "", err
	}

	return hash, err
}

// compare the password from the user with the database one
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("Failed to perform comparison: %v", err)
		return false, err
	}

	return match, err
}
