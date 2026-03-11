package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// hash the password using argon2id
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

// func to create and validate JWTs, which will be used to authenticate users
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    "chirpy-access",
		Subject:   userID.String(),
	})
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		// fmt.Printf("Error getting a token: %v", err)
		return "", err
	}

	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		// fmt.Printf("Error getting a token: %v", err)
		return uuid.Nil, err
	}

	subj, err := token.Claims.GetSubject()
	if err != nil {
		// fmt.Printf("Error getting subject")
		return uuid.Nil, err
	}

	uuidSubj, err := uuid.Parse(subj)
	if err != nil {
		fmt.Printf("Error parsing subject into UUID")
		return uuid.Nil, err
	}

	return uuidSubj, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		err := errors.New("Empty Auth Header")
		return "", err
	}

	splitHeader := strings.Fields(authHeader)
	if len(splitHeader) < 2 {
		err := errors.New("Wrong Auth Header")
		return "", err
	}
	return splitHeader[1], nil
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return ""
	}
	encodedStr := hex.EncodeToString(token)

	return encodedStr
}
