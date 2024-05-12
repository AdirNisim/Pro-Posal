package integrationtests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/pro-posal/webserver/api"
	"github.com/pro-posal/webserver/config"
	"github.com/pro-posal/webserver/internal/database"
	"github.com/pro-posal/webserver/services"
)

const URL = "http://localhost:8085"
const ADMIN_EMAIL = "admin@proposal.com"
const ADMIN_PASSWORD = "MyFancyAdminPassword"

var client *ApiClient

func TestMain(m *testing.M) {
	// TODO: Build and start a test database after running the migrations
	// These values will be changed to the fake database you'll start in the container
	config.AppConfig.Database.User = "admin"
	config.AppConfig.Database.Password = "Aa123456"
	config.AppConfig.Database.Name = "pro-posal"
	config.AppConfig.Database.Host = "localhost"
	config.AppConfig.Database.Port = "5432"
	config.AppConfig.Database.SSLMode = "disable"
	config.AppConfig.Server.Port = "8085"

	db := database.Connect()
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
	_, err := ums.CreateUser(context.Background(), services.CreateUserRequest{
		FirstName: "Joe",
		LastName:  "Doe",
		Email:     ADMIN_EMAIL,
		Password:  ADMIN_PASSWORD,
	})
	if err != nil && !strings.Contains(err.Error(), "users_email_hash_unique") {
		log.Fatalf("Failed seeding admin user: %v", err)
	}

	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)

	// Start the server
	log.Printf("Starting server on %s", addr)
	go func() {
		if err := http.ListenAndServe(addr, server.NewRouter()); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	client, err = NewApiClient(URL, ADMIN_EMAIL, ADMIN_PASSWORD)
	if err != nil {
		log.Fatalf("Failed creating API client: %v", err)
	}

	os.Exit(m.Run())
}
