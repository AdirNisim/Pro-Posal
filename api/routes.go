package api

import (
	"net/http"

	"github.com/pro-posal/webserver/internal/database"
	"github.com/pro-posal/webserver/internal/middlewares"
	"github.com/pro-posal/webserver/services"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type API struct {
	db                  *database.DBConnector
	userManagement      services.UserManagementService
	authService         services.AuthService
	companyManagement   services.CompanyManagementService
	premmisionManagment services.PermissionManagementService
	categoryManagment   services.CategoryManagementService
	contractManagment   services.ContractTemplateManagementService
	offerManagment      services.OfferManagementService
}

func NewAPI(
	db *database.DBConnector,
	userManagementService services.UserManagementService,
	authService services.AuthService,
	companyManagementService services.CompanyManagementService,
	premmisionManagment services.PermissionManagementService,
	categoryManagment services.CategoryManagementService,
	contractManagment services.ContractTemplateManagementService,
	offerManagment services.OfferManagementService,

) *API {
	return &API{
		db:                  db,
		userManagement:      userManagementService,
		authService:         authService,
		companyManagement:   companyManagementService,
		premmisionManagment: premmisionManagment,
		categoryManagment:   categoryManagment,
		contractManagment:   contractManagment,
		offerManagment:      offerManagment,
	}
}

func (a *API) NewRouter() http.Handler {
	router := mux.NewRouter()

	router.Use(middlewares.AccessLogMiddleware)
	router.Use(middlewares.PanicMiddleware)
	router.Use(middlewares.JSONHeaderMiddleware)
	router.Use(func(h http.Handler) http.Handler {
		return middlewares.AuthenticationMiddleware(a.authService, h)
	})

	// Check API status
	router.HandleFunc("/status", a.handleGetStatus).Methods("GET")

	// POST /users - Create a new user
	// GET /users - List all users
	router.HandleFunc("/users", a.PostUsers).Methods("POST")
	router.HandleFunc("/users", a.GetUsers).Methods("GET")

	// POST /users/login - Login user and obtain auth token
	router.HandleFunc("/users/login", a.PostUsersLogin).Methods("POST")
	// GET /users/{id} - Get user information
	router.HandleFunc("/users/{id}", a.GetUser).Methods("GET")
	// PATCH /users/updatePassword - update the users password
	router.HandleFunc("/users/updatePassword", a.UpdateUserPassword).Methods("PATCH")
	// Delete /users/{id} - Delete user
	router.HandleFunc("/users/{id}", a.DeleteUser).Methods("DELETE")

	// companies table
	// POST /companies - Create a new company
	router.HandleFunc("/companies", a.PostCompanies).Methods("POST")
	// GET /companies - User ID input List all companies
	router.HandleFunc("/companies/{id}", a.GetCompanies).Methods("GET")
	// PUT /companies/{id} - Company Id input Update company information
	router.HandleFunc("/companies/{id}", a.UpdateCompanies).Methods("PUT")
	// DELETE /companies/{id} - Company Id input Delete company
	router.HandleFunc("/companies/{id}", a.DeleteCompany).Methods("DELETE")

	// premmisions table
	// POST /premmisions/-> post a premmisions for company and email
	router.HandleFunc("/premmision/{id}", a.PostPermmision).Methods("POST")
	// GET //premmisions/{companyId} -> Get All Company users premmisions
	router.HandleFunc("/premmision/{id}", a.GetPermmisions).Methods("GET")
	// PUT /premmision/{id} -> Update premmision info
	router.HandleFunc("/premmision/{id}", a.UpdatePermission).Methods("PUT")
	// DELETE /premmision/{id} -> delete premmision
	router.HandleFunc("/premmision/{id}", a.DeletePermission).Methods("DELETE")

	// DELETE /companies/{companyID} -> Delete comapny from user

	// templates table
	// POST /contractsTemplates/{companyId} -> Post a  Compnay contract Template
	router.HandleFunc("/contractsTemplates/{companyId}", a.PostContractsTemplates).Methods("POST")
	// GET //contractsTemplates/{companyId} -> Get All Company contracts templates
	router.HandleFunc("/contractsTemplates/{companyId}", a.GetContractsTemplates).Methods("GET")
	// GET //contractsTemplates/{companyId}/{contractTemplateID} -> Get specific contract templates
	router.HandleFunc("/contractsTemplates/{id}", a.GetContractsTemplate).Methods("GET")
	// PUT /contractsTemplates/{companyId}/{contractTemplateID} -> Update contract template info
	router.HandleFunc("/contractsTemplates/{id}", a.UpdateContractsTemplates).Methods("PUT")
	// DELETE //contractsTemplates/{contractTemplateID} -> delete specic contract templates
	router.HandleFunc("/contractsTemplates/{id}", a.DeleteContractsTemplates).Methods("DELETE")

	// PUT /contractsTemplates/{companyId}/{contractTemplateID} -> Update contract template info
	// DELETE //contractsTemplates/{companyId}/{contractTemplateID} -> delete specic contract templates

	// categories table
	// POST /categories/{companyId} -> add a category for company
	// POST /categories/{companyId}/{categoryId} ->  add a sub category for company
	// GET /categories/{companyId} -> get categories for company
	// GET /categories/{companyId}/{categoryId}  -> Get sub categories of a category of a company
	// PUT /categories/{companyId}/{categoryId}/{ID} -> update sub category information
	// PUT /categories/{companyId}/{categoryID} -> update category id

	return router
}

// func handleJSON(w http.ResponseWriter, r *http.Request) {
// 	msg := map[string]string{"status": "ok"}
// 	response, err := json.Marshal(msg)

// 	if err != nil {
// 		log.Printf("failed marshaling response: %v", err)
// 		http.Error(w, "Something went wrong", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(response)
// }

// func handleRouteWithVariable(w http.ResponseWriter, r *http.Request) {
// 	// Use Gorilla Mux to extract the variable
// 	vars := mux.Vars(r)
// 	id := vars["id"]

// 	msg := map[string]string{"status": "ok", "id": id}
// 	response, err := json.Marshal(msg)

// 	if err != nil {
// 		log.Printf("failed marshaling response: %v", err)
// 		http.Error(w, "Something went wrong", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(response)
// }

// // handleUserPostRequest handles the POST request for creating a new user.
// // returns a 200 status code if successful.
// // Note: The handling of the 'invitedBy' input is pending implementation.

// func handleUserPostRequest(w http.ResponseWriter, r *http.Request) {
// 	var userInput userlib.UserInput
// 	err := json.NewDecoder(r.Body).Decode(&userInput)
// 	if err != nil {
// 		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
// 		return
// 	}

// 	passHash, err := userlib.HashPassword(userInput.Password)
// 	if err != nil {
// 		http.Error(w, "Error hashing password", http.StatusInternalServerError)
// 		return
// 	}

// 	if !userlib.IsValidEmail(userInput.Email) {
// 		http.Error(w, "Not a valid Email", http.StatusInternalServerError)
// 		return
// 	}

// 	// Need to know how to handle invitedBy input.
// 	newUser := models.User{
// 		ID:           uuid.New().String(),
// 		FirstName:    userInput.FirstName,
// 		LastName:     userInput.LastName,
// 		Phone:        userInput.Phone,
// 		Email:        userInput.Email,
// 		EmailHash:    userlib.HashEmail(userInput.Email),
// 		PasswordHash: passHash,
// 		InvitedBy:    null.String{},
// 		CreatedAt:    time.Now(),
// 		UpdatedAt:    time.Now(),
// 	}

// 	db, err := sql.Open("postgres", "user=admin password=Aa123456 dbname=pro-posal host=localhost sslmode=disable")
// 	if err != nil {
// 		log.Fatalf("Could not connect to the database! %v", err)
// 	}
// 	defer db.Close() // Will exec in the end of main

// 	ctx := r.Context()
// 	err = newUser.Insert(ctx, db, boil.Infer())
// 	if err != nil {
// 		log.Fatalf("Error inserting user: %v", err)
// 		http.Error(w, "User is already in the system. Please log-in!", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	log.Printf("User created successfully")
// }

// // handleUserGetRequest handles the GET request for user information retrieval.
// // If successful, it returns the user's information in JSON format with status code 200 (OK).
// // Otherwise, it sends appropriate error responses
// func handleUserGetRequest(w http.ResponseWriter, r *http.Request) {
// 	var userInput userlib.LoginRequest
// 	err := json.NewDecoder(r.Body).Decode(&userInput)
// 	if err != nil {
// 		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
// 		return
// 	}

// 	user, err := userlib.GetUser(userInput.Email, userInput.Password, r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	jsonData, err := json.Marshal(user)
// 	if err != nil {
// 		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
// 		log.Printf("Error marshaling JSON")
// 		return
// 	}

// 	log.Printf("User was successfully sent.")
// 	w.Write(jsonData)
// }

// // func handleUserUpdatePassword(w http.ResponseWriter, r *http.Request) {
// // 	var userInput userlib.LoginRequest
// // 	err := json.NewDecoder(r.Body).Decode(&userInput)
// // 	if err != nil {
// // 		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
// // 		return
// // 	}

// // 	user, err := userlib.GetUser(userInput.Email, userInput.Password, r)
// // 	if err != nil {
// // 		http.Error(w, err.Error(), http.StatusUnauthorized)
// // 		return
// // 	}

// // }
