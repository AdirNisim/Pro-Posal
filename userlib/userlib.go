package userlib

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"webserver/models"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func HashEmail(email string) string {
	hasher := sha256.New()
	hasher.Write([]byte(email))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateUserFromInput(input UserInput) (*models.User, error) {
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	hashedEmail := HashEmail(input.Email)

	return &models.User{
		ID:           uuid.New().String(),
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Phone:        input.Phone,
		Email:        input.Email,
		EmailHash:    hashedEmail,
		PasswordHash: hashedPassword,
		InvitedBy:    null.NewString("", false),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
