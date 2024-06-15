package integrationtests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pro-posal/webserver/api"
	"github.com/pro-posal/webserver/config"
	"github.com/pro-posal/webserver/internal/database"
	"github.com/pro-posal/webserver/services"
)

var client *ApiClient

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Step 1: Run a new Docker container
	host, port, postgres, err := database.RunDockerContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to run Docker container: %v", err)
	}
	defer postgres.Terminate(ctx)

	// Update the configuration to use the container's host and port
	config.TestConfig.TestDatabase.Host = host
	config.TestConfig.TestDatabase.Port = port

	// Step 2: Run migrations
	if err := database.RunMigrations(); err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return
	}

	log.Println("Database setup completed successfully")

	db := database.TestConnect()
	defer db.Conn.Close()

	ums := services.NewUserManagementService(db)
	auth := services.NewAuthService(db)
	cms := services.NewCompanyManagementService(db)
	pms := services.NewPermissionManagementService(db)
	cams := services.NewCategoryManagementService(db)
	ctms := services.NewContractTemplateManagementService(db)
	oms := services.NewOfferManagementService(db)

	server := api.NewAPI(db, ums, auth, cms, pms, cams, ctms, oms)

	// Seed an admin user
	_, err = ums.CreateUser(context.Background(), services.CreateUserRequest{
		FirstName: config.TestConfig.User.FirstName,
		LastName:  config.TestConfig.User.LastName,
		Phone:     config.TestConfig.User.Phone,
		Email:     config.TestConfig.User.Email,
		Password:  config.TestConfig.User.Password,
	})
	if err != nil && !strings.Contains(err.Error(), "users_email_hash_unique") {
		log.Fatalf("Failed seeding admin user: %v", err)
	}

	addr := fmt.Sprintf(":%s", config.TestConfig.TestServer.Port)

	// Start the server
	log.Printf("Starting server on %s", addr)
	go func() {
		fmt.Println("Starting http server...")
		if err := http.ListenAndServe(addr, server.NewRouter()); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// // Wait for the server to start
	serverAddress := config.TestConfig.TestServer.URL
	if err := waitForService(serverAddress, 30*time.Second); err != nil {
		log.Fatalf("Failed to wait for server: %v", err)
	}

	client, err = NewApiClient(serverAddress,
		config.TestConfig.User.Email,
		config.TestConfig.User.Password)
	if err != nil {
		log.Fatalf("Failed creating API client: %v", err)
	}

	os.Exit(m.Run())
}

func waitForService(serviceURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		_, err := http.Get(serviceURL)
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("service at %s is not available", serviceURL)
}
