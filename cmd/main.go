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

/*		companyManagement:   companyManagementService,
		premmisionManagment: premmisionManagment,
		categoryManagment:   categoryManagment,
		contractManagment:   contractManagment,
		offerManagment:      offerManagment,*/

func main() {
	db := database.Connect()
	defer db.Conn.Close()

	// Initialize router (for different api routes)
	ums := services.NewUserManagementService(db)
	auth := services.NewAuthService(db)
	cms := services.NewCompanyManagementService(db)
	pms := services.NewPermissionManagementService(db)
	cams := services.NewCategoryManagementService(db)
	ctms := services.NewContractTemplateManagementService(db)
	oms := services.NewOfferManagementService(db)

	server := api.NewAPI(db, ums, auth, cms, pms, cams, ctms, oms)

	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)

	// // Start the server
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, server.NewRouter()); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
