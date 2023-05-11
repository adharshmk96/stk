package utils

import (
	"crypto/rand"
	"encoding/base64"
	"reflect"

	"golang.org/x/crypto/argon2"
)

// Usage
/*
salt, err := generateSalt()
if err != nil {
	log.Fatal("Failed to generate salt:", err)
}
*/
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// Usage
// NOTE: Store both the hash and salt in the database, it is needed to verify the password
/*
hash, salt := hashPassword(password, salt)
*/
func HashPassword(password string, salt []byte) (string, string) {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	hashedPassword := base64.RawStdEncoding.EncodeToString(hash)
	hashedSalt := base64.RawStdEncoding.EncodeToString(salt)

	return hashedPassword, hashedSalt
}

// Usage
/*
isValid, err := verifyPassword(hashedPassword, hashedSalt, password)
if err != nil {
	log.Fatal("Failed to verify password:", err)
}
*/
func VerifyPassword(storedHash, storedSalt, password string) (bool, error) {
	salt, err := base64.RawStdEncoding.DecodeString(storedSalt)
	if err != nil {
		return false, err
	}

	newHash, _ := HashPassword(password, salt)
	return reflect.DeepEqual(storedHash, newHash), nil
}
