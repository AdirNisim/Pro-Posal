package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashEmail(email string) string {
	lowerEmail := strings.ToLower(email)
	hasher := sha256.New()
	hasher.Write([]byte(lowerEmail))
	return hex.EncodeToString(hasher.Sum(nil))
}
