package main

//Some change
import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleJSON)
	http.ListenAndServe(":8080", nil)
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	msg := map[string]string{"status": "ok"}
	respons, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respons)
}
