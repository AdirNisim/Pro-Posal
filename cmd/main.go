package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pro-posal/webserver/api"
	"github.com/pro-posal/webserver/config"
	"github.com/pro-posal/webserver/internal/database"
	"github.com/pro-posal/webserver/services"
)

func main() {
	db := database.Connect()
	defer db.Conn.Close()

	// Initialize router (for different api routes)
	ums := services.NewUserManagementService(db)
	auth := services.NewAuthService(db)
	server := api.NewAPI(db, ums, auth)

	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)

	// // Start the server
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, server.NewRouter()); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
