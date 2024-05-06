package api

import (
	"encoding/json"
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

type PUTPremmisionRequestBody struct {
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

	company, err := a.premmisionManagment.CreatePermission(r.Context(), services.CreatePermissionRequest{
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

	resp, err := json.Marshal(company)
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling Company response object", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusCreated)

}

func (a *API) GetPermmisions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["id"]
	usersPermissions, err := a.premmisionManagment.GetPermissions(r.Context(), companyId)
	if err != nil {
		log.Printf("Error getting permissions: %v", err)
		http.Error(w, "Error getting permissions", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(GetUsersPermissionsResponseBody{TotalUsers: len(usersPermissions), UserPermissions: usersPermissions})
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling companies response object", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}

func (a *API) DeletePermmision(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	permissionId := vars["id"]

	permission, err := a.premmisionManagment.DeletePermission(r.Context(), permissionId)
	if err != nil {
		log.Printf("Error deleting permission: %v", err)
		http.Error(w, "Error deleting permission", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(permission)
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling Company response object", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusCreated)
}

func (a *API) UpdatePermmision(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	permissionId := vars["id"]
	var request PUTPremmisionRequestBody

	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	permission, err := a.premmisionManagment.UpdatePermission(r.Context(), services.UpdatePermissionRequest{
		Id:         permissionId,
		Role:       request.Role,
		ContractID: request.ContractID,
	})

	if err != nil {
		log.Printf("Error updating permission: %v", err)
		http.Error(w, "Error updating permission", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(permission)
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling Company response object", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusCreated)
}
