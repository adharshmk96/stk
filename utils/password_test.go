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

	t.Run("wrong password", func(t *testing.T) {
		invalidPassword := "wrongpassword"
		isValid, err := VerifyPassword(hashedPassword, hashedSalt, invalidPassword)

		if err != nil {
			t.Fatalf("VerifyPassword() error with invalid password: %v", err)
		}

		if isValid {
			t.Fatalf("VerifyPassword() validated an incorrect password")
		}
	})

	t.Run("Invalid hashed salt", func(t *testing.T) {
		invalidHashedSalt := "invalidhashedsalt"
		isValid, err := VerifyPassword(hashedPassword, invalidHashedSalt, password)
		if err == nil {
			t.Fatalf("VerifyPassword() did not return an error with an invalid hashed salt")
		}

		if isValid {
			t.Fatalf("VerifyPassword() validated an invalid hashed salt")
		}
	})

	t.Run("Invalid hashed password", func(t *testing.T) {
		invalidHashedPassword := "invalidhashedpassword"
		isValid, err := VerifyPassword(invalidHashedPassword, hashedSalt, password)
		if err != nil {
			t.Fatalf("VerifyPassword() error with a random password hash string: %v", err)
		}

		if isValid {
			t.Fatalf("VerifyPassword() validated an invalid hashed password")
		}
	})

	t.Run("valid password", func(t *testing.T) {
		isValid, err := VerifyPassword(hashedPassword, hashedSalt, password)
		if err != nil {
			t.Fatalf("VerifyPassword() error: %v", err)
		}

		if !isValid {
			t.Fatalf("VerifyPassword() failed to verify the correct password")
		}
	})

}
