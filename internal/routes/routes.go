package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"webserver/models"
	"webserver/userlib"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func NewRouter() http.Handler {

	router := mux.NewRouter()
	router.Use(middleware)
	router.HandleFunc("/", handleJSON).Methods("GET")
	router.HandleFunc("/api/templates/{id}", handleRouteWithVariable).Methods("GET")
	router.HandleFunc("/v1/addUser", handleUserPostRequest).Methods("POST")
	router.HandleFunc("/v1/getUser", handleUserGetRequest).Methods("POST")
	return router
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	msg := map[string]string{"status": "ok"}
	response, err := json.Marshal(msg)

	if err != nil {
		log.Printf("failed marshaling response: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func handleRouteWithVariable(w http.ResponseWriter, r *http.Request) {
	// Use Gorilla Mux to extract the variable
	vars := mux.Vars(r)
	id := vars["id"]

	msg := map[string]string{"status": "ok", "id": id}
	response, err := json.Marshal(msg)

	if err != nil {
		log.Printf("failed marshaling response: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// handleUserPostRequest handles the POST request for creating a new user.
// returns a 200 status code if successful.
// Note: The handling of the 'invitedBy' input is pending implementation.

func handleUserPostRequest(w http.ResponseWriter, r *http.Request) {
	var userInput userlib.UserInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	passHash, err := userlib.HashPassword(userInput.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	if !userlib.IsValidEmail(userInput.Email) {
		http.Error(w, "Not a valid Email", http.StatusInternalServerError)
		return
	}

	// Need to know how to handle invitedBy input.
	newUser := models.User{
		ID:           uuid.New().String(),
		FirstName:    userInput.FirstName,
		LastName:     userInput.LastName,
		Phone:        userInput.Phone,
		Email:        userInput.Email,
		EmailHash:    userlib.HashEmail(userInput.Email),
		PasswordHash: passHash,
		InvitedBy:    null.String{},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	db, err := sql.Open("postgres", "user=admin password=Aa123456 dbname=pro-posal host=localhost sslmode=disable")
	if err != nil {
		log.Fatalf("Could not connect to the database! %v", err)
	}
	defer db.Close() // Will exec in the end of main

	ctx := r.Context()
	err = newUser.Insert(ctx, db, boil.Infer())
	if err != nil {
		log.Fatalf("Error inserting user: %v", err)
		http.Error(w, "User is already in the system. Please log-in!", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("User created successfully")
}

// handleUserGetRequest handles the GET request for user information retrieval.
// If successful, it returns the user's information in JSON format with status code 200 (OK).
// Otherwise, it sends appropriate error responses
func handleUserGetRequest(w http.ResponseWriter, r *http.Request) {
	var userInput userlib.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	user, err := userlib.GetUser(userInput.Email, userInput.Password, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		log.Printf("Error marshaling JSON")
		return
	}

	log.Printf("User was successfully sent.")
	w.Write(jsonData)
}

// func handleUserUpdatePassword(w http.ResponseWriter, r *http.Request) {
// 	var userInput userlib.LoginRequest
// 	err := json.NewDecoder(r.Body).Decode(&userInput)
// 	if err != nil {
// 		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
// 		return
// 	}

// 	user, err := userlib.GetUser(userInput.Email, userInput.Password, r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// }
