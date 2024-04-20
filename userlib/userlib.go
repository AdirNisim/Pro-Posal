package userlib

// import (
// 	"crypto/sha256"
// 	"database/sql"
// 	"encoding/hex"
// 	"errors"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/pro-posal/webserver/models"

// 	"github.com/asaskevich/govalidator"
// 	"github.com/volatiletech/sqlboiler/v4/queries/qm"
// 	"golang.org/x/crypto/bcrypt"
// )

// type UserInput struct {
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Phone     string `json:"phone"`
// 	Email     string `json:"email"`
// 	Password  string `json:"password"`
// }

// type LoginRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// // HashPassword securely hashes the given plaintext password using bcrypt encryption
// // with a specified cost factor (12 will take 250 milsec and double with each increment)
// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
// 	return string(bytes), err
// }

// // CheckPassword compares a hashed password with a plaintext password
// // to verify whether they match. It uses the bcrypt package to securely
// // compare the hashed password with the plaintext password.
// func CheckPassword(hashedPassword, password string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
// 	return err == nil
// }

// // IsValidEmail checks whether the given string is a valid email address
// // using the govalidator package. It performs basic validation to ensure
// // that the email address follows standard formatting rules, including
// // having a valid structure (e.g., user@example.com).
// func IsValidEmail(email string) bool {
// 	return govalidator.IsEmail(email)
// }

// // HashEmail hashes the given email address using SHA-256 encryption
// // while ensuring all characters are converted to lowercase for consistency.
// func HashEmail(email string) string {
// 	lowerEmail := strings.ToLower(email)
// 	hasher := sha256.New()
// 	hasher.Write([]byte(lowerEmail))
// 	return hex.EncodeToString(hasher.Sum(nil))
// }

// // GetUser retrieves a user from the database based on the provided email hash.
// // It attempts to authenticate the user's password with the database and returns
// // a pointer to the user's struct if authentication is successful and error otherwise.
// func GetUser(email string, password string, r *http.Request) (*models.User, error) {
// 	db, err := sql.Open("postgres", "user=admin password=Aa123456 dbname=pro-posal host=localhost sslmode=disable")
// 	if err != nil {
// 		return nil, errors.New("could not connect to the database")
// 	}
// 	defer db.Close() // Will exec in the end of main

// 	ctx := r.Context()
// 	emailHash := HashEmail(email)
// 	user, err := models.Users(qm.Where("email_hash=?", emailHash)).One(ctx, db)
// 	if err != nil {
// 		log.Printf("No user found with the provided email: %s", emailHash)
// 		return nil, errors.New("password or email does not match our records")
// 	}

// 	if !CheckPassword(user.PasswordHash, password) {
// 		log.Printf("User authentication failed!")
// 		return nil, errors.New("password or email does not match our records")
// 	}
// 	return user, nil
// }
