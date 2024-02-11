package utils

import "golang.org/x/crypto/bcrypt"

// helper function that uses bcrypt to hash a plain password
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

// check that plain password provides the same hash
func CheckPassword(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
