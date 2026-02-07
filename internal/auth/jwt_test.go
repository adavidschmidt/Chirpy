package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestJWT_isAccepted(t *testing.T) {
	userID := uuid.New()
	userSecret := "secretsecret"

	tokenString, err := MakeJWT(userID, userSecret)
	if err != nil {
		t.Fatalf("Error creating the token: %v", err)
	}
	userJWTID, err := ValidateJWT(tokenString, userSecret)
	if err != nil {
		t.Fatalf("Error validating token: %v", err)
	}
	if userID != userJWTID {
		t.Fatalf("Returned userID does not equal the provided userID. Got: %v, want: %v", userJWTID, userID)
	}
}

func TestExpiredToken_isRejected(t *testing.T) {
	userID := uuid.New()
	userSecret := "secretsecret"

	tokenString, err := MakeJWT(userID, userSecret)
	if err != nil {
		t.Fatalf("Error creating the token: %v", err)
	}
	userJWTID, err := ValidateJWT(tokenString, userSecret)
	if err == nil {
		t.Fatal("Token is expired, but was validated")
	}
	if userJWTID != uuid.Nil {
		t.Fatalf("expected uuid.Nil on error, got %v", userJWTID)
	}
}

func TestWrongSecret_isRejected(t *testing.T) {
	userID := uuid.New()
	userSecret := "secretsecret"

	tokenString, err := MakeJWT(userID, userSecret)
	if err != nil {
		t.Fatalf("Error creating the token: %v", err)
	}
	wrongSecret := "wrongsecret"
	userJWTID, err := ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Fatal("Wrong secret was provided but no error triggered when validating token")
	}
	if userJWTID != uuid.Nil {
		t.Fatalf("expected uuid.Nil on error, got %v", userJWTID)
	}
}
