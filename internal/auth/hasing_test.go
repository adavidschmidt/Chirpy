package auth

import (
	"testing"

	"github.com/alexedwards/argon2id"
)

// Tests written by Copilot
func TestHashPassword(t *testing.T) {
	password := "supersecret"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == "" {
		t.Fatalf("expected non-empty hash")
	}

	// Verify the hash is valid Argon2id
	ok, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		t.Fatalf("argon2id.ComparePasswordAndHash returned error: %v", err)
	}
	if !ok {
		t.Fatalf("expected hash to match password")
	}
}

func TestCheckPasswordHash_Success(t *testing.T) {
	password := "supersecret"

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Fatalf("failed to create hash: %v", err)
	}

	ok, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if !ok {
		t.Fatalf("expected CheckPasswordHash to return true for valid password")
	}
}

func TestCheckPasswordHash_Failure(t *testing.T) {
	password := "supersecret"
	wrongPassword := "nottherightone"

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Fatalf("failed to create hash: %v", err)
	}

	ok, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}

	if ok {
		t.Fatalf("expected CheckPasswordHash to return false for invalid password")
	}
}
