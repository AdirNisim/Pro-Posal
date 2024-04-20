package main

import (
	"fmt"
	"log"
	"net/http"
	"webserver/internal/routes"
)

func main() {
	err := godotenv.Load("../_.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	
	// Initialize router (for different api routes)
	router := routes.NewRouter()

	addr := fmt.Sprintf(":%s", AppConfig.Port)

	// Start the server
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
