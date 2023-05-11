package utils

import (
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error: %v", err)
	}
	if len(salt) != 16 {
		t.Fatalf("GenerateSalt() returned incorrect salt length: got %v, want 16", len(salt))
	}
}

func TestHashPassword(t *testing.T) {
	password := "supersecretpassword"
	salt, _ := GenerateSalt()

	hashedPassword, hashedSalt := HashPassword(password, salt)

	if hashedPassword == "" {
		t.Fatalf("HashPassword() returned an empty hashed password")
	}

	if hashedSalt == "" {
		t.Fatalf("HashPassword() returned an empty hashed salt")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "supersecretpassword"
	salt, _ := GenerateSalt()

	hashedPassword, hashedSalt := HashPassword(password, salt)

	isValid, err := VerifyPassword(hashedPassword, hashedSalt, password)
	if err != nil {
		t.Fatalf("VerifyPassword() error: %v", err)
	}

	if !isValid {
		t.Fatalf("VerifyPassword() failed to verify the correct password")
	}

	invalidPassword := "wrongpassword"
	isValid, err = VerifyPassword(hashedPassword, hashedSalt, invalidPassword)

	if err != nil {
		t.Fatalf("VerifyPassword() error with invalid password: %v", err)
	}

	if isValid {
		t.Fatalf("VerifyPassword() validated an incorrect password")
	}
}
