package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pro-posal/webserver/internal/utils"
	"github.com/pro-posal/webserver/models"
	"github.com/pro-posal/webserver/services"
)

type PostUsersRequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type PutUsersPasswordRequestBody struct {
	NewPassword string `json:"new_password"`
}

func (a *API) PostUsers(w http.ResponseWriter, r *http.Request) {
	var request PostUsersRequestBody
	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// TODO: Add some input validations

	user, err := a.userManagement.CreateUser(r.Context(), services.CreateUserRequest{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Phone:     request.Phone,
		Email:     strings.ToLower(request.Email),
		Password:  request.Password,
	})

	if err != nil {
		log.Printf("Failed creating user: %v", err)
		http.Error(w, "Failed creating user", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(user)
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling user response object", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusCreated)
}

type GetUsersResponseBody struct {
	TotalUsers int            `json:"total_users"`
	Users      []*models.User `json:"users"`
}

func (a *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.userManagement.ListUsers(r.Context())
	if err != nil {
		log.Printf("Failed listing users: %v", err)
		http.Error(w, "Failed listing users", http.StatusInternalServerError)
		return
	}

	// Example - how do fetch the session and filter by context
	// session := r.Context().Value("session").(*models.Session)
	// log.Printf("Request is invoked by user %v", session.UserID)

	// var filtered []*models.User
	// for _, user := range users {
	// 	if user.ID == session.UserID.String() {
	// 		filtered = append(filtered, user)
	// 	}
	// }

	resp, err := json.Marshal(GetUsersResponseBody{TotalUsers: len(users), Users: users})
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling users response object", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}

func (a *API) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := a.userManagement.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("Failed retrieving user: %v", err)
		http.Error(w, "Failed retrieving user", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(user)
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling user response object", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (a *API) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	var request PutUsersPasswordRequestBody
	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	user, err := a.userManagement.UpdateUserPassword(r.Context(), services.ChangeUserPasswordRequest{
		Uuid:        "SOMEUUID",
		NewPassword: request.NewPassword,
	})
	if err != nil {
		log.Printf("Error updating user password: %v", err)
		http.Error(w, "Error updating user password", http.StatusInternalServerError)
		return
	}

	// Assuming there's a method to serialize the user back to JSON or some response
	response, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		http.Error(w, "Error preparing response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
