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

type PostCompanyRequestBody struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	ContactID  string `json:"contact_id"`
	LogoBase64 string `json:"logo_base64"`
}
type PUTCompanyRequestBody struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	LogoBase64 string `json:"logo_base64"`
}

type GetCompaniesResponseBody struct {
	TotalCompanies int               `json:"total_companies"`
	Companies      []*models.Company `json:"companies"`
}

func (a *API) PostCompanies(w http.ResponseWriter, r *http.Request) {
	var request PostCompanyRequestBody
	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}
	company, err := a.companyManagement.CreateCompany(r.Context(), services.CreateCompanyRequest{
		Name:       request.Name,
		ContactID:  request.ContactID,
		Address:    request.Address,
		LogoBase64: request.LogoBase64})
	if err != nil {
		log.Printf("Error Create a Comapny / Adding a user for company: %v", err)
		http.Error(w, "Error Create a Comapny / Adding a user for company", http.StatusBadRequest)
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

func (a *API) GetCompanies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contactId := vars["id"]

	companies, err := a.companyManagement.GetCompanies(r.Context(), contactId)
	if err != nil {
		log.Printf("Error Getting Companies for this ID: %v", err)
		http.Error(w, "Error Getting Companies", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(GetCompaniesResponseBody{TotalCompanies: len(companies), Companies: companies})
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling companies response object", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}

func (a *API) UpdateCompanies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["id"]
	var request PUTCompanyRequestBody

	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	company, err := a.companyManagement.UpdateCompany(r.Context(), companyId, services.UpdateCompanyRequest{
		Name:       request.Name,
		Address:    request.Address,
		LogoBase64: request.LogoBase64,
	})

	if err != nil {
		log.Printf("Error Create a Comapny / Adding a user for company: %v", err)
		http.Error(w, "Error Create a Comapny / Adding a user for company", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(company)
	if err != nil {
		log.Printf("Failed marshaling response: %v", err)
		http.Error(w, "Failed marshalling Company response object", http.StatusInternalServerError)
		return
	}

	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}
