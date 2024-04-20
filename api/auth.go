package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/pro-posal/webserver/internal/utils"
)

type PostUsersLoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PostUsersLoginResponseBody struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (a *API) PostUsersLogin(w http.ResponseWriter, r *http.Request) {
	var request PostUsersLoginRequestBody
	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	token, expires, err := a.authService.CreateAuthToken(r.Context(), strings.ToLower(request.Email), request.Password)
	if err != nil {
		if strings.Contains(err.Error(), "invalid email or password") {
			http.Error(w, "Invalid email or password", http.StatusForbidden)
			return
		}

		log.Printf("Failed creating auth token: %v", err)
		http.Error(w, "Failed creating auth token", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(PostUsersLoginResponseBody{
		AccessToken: token,
		ExpiresAt:   expires.UnixMilli(),
	})
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed creating response object", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}
