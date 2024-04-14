package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"webserver/internal/routes"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../_.env") // This will look for a ".env" file in the current directory
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Initialize router (for different api routes)
	router := routes.NewRouter()

	// Read port from env file
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)

	// Start the server
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
