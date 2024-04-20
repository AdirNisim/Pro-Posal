package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

	resp, err := json.Marshal(GetUsersResponseBody{TotalUsers: len(users), Users: users})
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling users response object", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}
