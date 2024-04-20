package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a *API) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	err := a.db.Conn.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		http.Error(w, "Failed to communicate with database", http.StatusInternalServerError)
		return
	}

	msg := map[string]string{"status": "ok"}
	response, err := json.Marshal(msg)

	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
