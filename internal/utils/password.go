package utils

import "golang.org/x/crypto/bcrypt"

// GeneratePassword generates a bcrypt hash for the given password.
func GeneratePassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// ComparePassword verifies if the given password matches the stored hash.
func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
