package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {

	router := mux.NewRouter()
	router.Use(middleware)
	router.HandleFunc("/", handleJSON).Methods("GET")
	router.HandleFunc("/api/templates/{id}", handleRouteWithVariable).Methods("GET")

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
