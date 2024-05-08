package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/pro-posal/webserver/models"
)

func UnmarshalRequest(r *http.Request, v any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return fmt.Errorf("failed decoding request JSON: %w", err)
	}

	return nil
}

func MarshalAndWriteResponse(w http.ResponseWriter, data interface{}) {
	resp, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling response object", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func GetUserIDFromSession(r *http.Request) uuid.UUID {
	session, ok := r.Context().Value("session").(*models.Session)
	if !ok {
		log.Println("Session is not found or is of the incorrect type")
		return uuid.Nil
	}

	log.Printf("Request is invoked by user %v", session.UserID)
	return uuid.UUID(session.UserID)
}
