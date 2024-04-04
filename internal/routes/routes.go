package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {

	router := mux.NewRouter()

	router.HandleFunc("/", handleJSON).Methods("GET")
	router.HandleFunc("/api/templates/{id}", handleRouteWithVariable).Methods("GET")

	return router
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	msg := map[string]string{"status": "ok"}
	response, err := json.Marshal(msg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func handleRouteWithVariable(w http.ResponseWriter, r *http.Request) {

	// Use Gorilla Mux to extract the variable
	vars := mux.Vars(r)
	id := vars["id"]

	msg := map[string]string{"status": "ok", "id": id}
	response, err := json.Marshal(msg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
