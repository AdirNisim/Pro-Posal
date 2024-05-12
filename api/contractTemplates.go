package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pro-posal/webserver/internal/utils"
	"github.com/pro-posal/webserver/models"
	"github.com/pro-posal/webserver/services"
)

type PostContractRequestBody struct {
	Name      string `json:"name"`
	Template  string `json:"template"`
	CompanyID string `json:"company_id"`
}

type PUTContractRequestBody struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

type GetContractResponseBody struct {
	TotalCompanies   int                        `json:"total_contract_templates"`
	ContractTemplate []*models.ContractTemplate `json:"contract_templates"`
}

func (a *API) PostContractsTemplates(w http.ResponseWriter, r *http.Request) {
	var request PostContractRequestBody
	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}
	contract, err := a.contractManagment.PostContractsTemplate(r.Context(), services.CreateContractTemplateRequest{
		Name:      request.Name,
		Template:  request.Template,
		CompanyID: request.CompanyID})
	if err != nil {
		log.Printf("Error Creating a Contract Template: %v", err)
		http.Error(w, "Error Creating a Contract Template", http.StatusBadRequest)
		return
	}

	utils.MarshalAndWriteResponse(w, contract)
}

func (a *API) GetContractsTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	contracts, err := a.contractManagment.GetContractsTemplate(r.Context(), id)
	if err != nil {
		log.Printf("Error Getting Contract Template: %v", err)
		http.Error(w, "Error Getting Contract Template", http.StatusBadRequest)
		return
	}
	utils.MarshalAndWriteResponse(w, contracts)
}

func (a *API) GetContractsTemplates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID := vars["companyId"]

	contracts, err := a.contractManagment.GetContractsTemplates(r.Context(), companyID)
	if err != nil {
		log.Printf("Error Getting Contract Templates: %v", err)
		http.Error(w, "Error Getting Contract Templates", http.StatusBadRequest)
		return
	}
	responseBody := GetContractResponseBody{TotalCompanies: len(contracts), ContractTemplate: contracts}
	utils.MarshalAndWriteResponse(w, responseBody)
}

func (a *API) UpdateContractsTemplates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var request PUTContractRequestBody

	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}
	contract, err := a.contractManagment.UpdateContractsTemplate(r.Context(), id, services.UpdateContractsTemplatesRequest{
		Name:     request.Name,
		Template: request.Template})
	if err != nil {
		log.Printf("Error Updating Contract Template: %v", err)
		http.Error(w, "Error Updating Contract Template", http.StatusBadRequest)
		return
	}
	utils.MarshalAndWriteResponse(w, contract)
}

func (a *API) DeleteContractsTemplates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	contract, err := a.contractManagment.DeleteContractsTemplate(r.Context(), id)
	if err != nil {
		log.Printf("Error Deleting Contract Template: %v", err)
		http.Error(w, "Error Deleting Contract Template", http.StatusBadRequest)
		return
	}
	utils.MarshalAndWriteResponse(w, contract)
}
