package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pro-posal/webserver/internal/utils"
	"github.com/pro-posal/webserver/models"
	"github.com/pro-posal/webserver/services"
)

type PostPremmisionRequestBody struct {
	UserID     string `json:"user_id"`
	CompanyID  string `json:"company_id"`
	Role       string `json:"role"`
	ContractID string `json:"contract_id"`
}

type UpdatePremmisionRequestBody struct {
	Role       string `json:"role"`
	ContractID string `json:"contract_id"`
}

type GetUsersPermissionsResponseBody struct {
	TotalUsers      int                      `json:"total_users"`
	UserPermissions []*models.UserPermission `json:"users_permissions"`
}

func (a *API) PostPermmision(w http.ResponseWriter, r *http.Request) {
	var request PostPremmisionRequestBody
	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}
	// Check that request.Role is one of the predefined roles
	validRole := request.Role == string(models.CompanyAdminRole) ||
		request.Role == string(models.CompanyContributorRole) ||
		request.Role == string(models.CompanyProjectManagerRole) ||
		request.Role == string(models.ProspectRole)

	if !validRole {
		log.Printf("Invalid role provided: %v", request.Role)
		http.Error(w, "Invalid role provided", http.StatusBadRequest)
		return
	}

	permission, err := a.permissionsManagement.CreatePermission(r.Context(), services.CreatePermissionRequest{
		UserID:     request.UserID,
		CompanyID:  request.CompanyID,
		Role:       request.Role,
		ContractID: request.ContractID,
	})
	if err != nil {
		log.Printf("Error creating a user's permission for company: %v", err)
		http.Error(w, "Error creating a user's permission for company", http.StatusBadRequest)
		return
	}

	utils.MarshalAndWriteResponse(w, permission)

}

func (a *API) GetPermmisions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["id"]
	usersPermissions, err := a.permissionsManagement.GetPermissions(r.Context(), companyId)
	if err != nil {
		log.Printf("Error getting permissions: %v", err)
		http.Error(w, "Error getting permissions", http.StatusBadRequest)
		return
	}

	responseBody := GetUsersPermissionsResponseBody{
		TotalUsers:      len(usersPermissions),
		UserPermissions: usersPermissions,
	}

	utils.MarshalAndWriteResponse(w, responseBody)
}

func (a *API) DeletePermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	permissionId := vars["id"]

	permission, err := a.permissionsManagement.DeletePermission(r.Context(), permissionId)
	if err != nil {
		log.Printf("Error deleting permission: %v", err)
		http.Error(w, "Error deleting permission", http.StatusBadRequest)
		return
	}

	utils.MarshalAndWriteResponse(w, permission)
}

func (a *API) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	permissionId := vars["id"]

	var request UpdatePremmisionRequestBody
	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	permission, err := a.permissionsManagement.UpdatePermission(r.Context(), services.UpdatePermissionRequest{
		Id:         permissionId,
		Role:       request.Role,
		ContractID: request.ContractID,
	})
	if err != nil {
		log.Printf("Error updating a user's permission for company: %v", err)
		http.Error(w, "Error updating a user's permission for company", http.StatusBadRequest)
		return
	}

	utils.MarshalAndWriteResponse(w, permission)
}
