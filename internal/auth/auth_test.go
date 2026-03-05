package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	u := uuid.New()
	tokenSecret := "abcde"
	duration := time.Duration(50) * time.Second
	signedString, err := MakeJWT(u, tokenSecret, duration)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	u2, err := ValidateJWT(signedString, tokenSecret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	if u2 != u {
		t.Errorf("expected %v, got %v", u, u2)
	}
}

func TestWrongKey(t *testing.T) {
	u := uuid.New()
	tokenSecret := "abcde"
	tokenSecretWrong := "edcba"
	duration := time.Duration(50) * time.Second
	signedString, err := MakeJWT(u, tokenSecret, duration)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	u2, err := ValidateJWT(signedString, tokenSecretWrong)

	if err == nil {
		t.Errorf("expected an error due to wrong secret, but got none")
	}

	if u2 != uuid.Nil {
		t.Errorf("expected uuid.Nil, got %v", u2)
	}
}

func TestExpiredKey(t *testing.T) {
	u := uuid.New()
	tokenSecret := "abcde"
	duration := time.Duration(-50) * time.Second
	signedString, err := MakeJWT(u, tokenSecret, duration)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}

	u2, err := ValidateJWT(signedString, tokenSecret)

	if err == nil {
		t.Errorf("expected an error due to expiration")
	}

	if u2 != uuid.Nil {
		t.Errorf("expected uuid.Nil, got %v", u2)
	}
}
