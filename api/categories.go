package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pro-posal/webserver/internal/utils"
	"github.com/pro-posal/webserver/models"
	"github.com/pro-posal/webserver/services"
)

type PostCategoriesRequestBody struct {
	CompanyID   string `json:"company_id"`
	CategoryID  string `json:"category_id"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
type PutCategoriesRequestBody struct {
	CategoryID  string `json:"category_id"`
	Description string `json:"description"`
}

type GetCategoriesResponseBody struct {
	TotalCompanies int                `json:"total_companies"`
	Description    []*models.Category `json:"categories"`
}

func (a *API) PostCategories(w http.ResponseWriter, r *http.Request) {
	var request PostCategoriesRequestBody
	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}
	category, err := a.categoryManagment.CreateCategory(r.Context(), services.CreateCategoryRequest{
		CompanyID:   request.CompanyID,
		Description: request.Description,
		Type:        request.Type})
	if err != nil {
		log.Printf("Error Creating a Category: %v", err)
		http.Error(w, "Error Creating a Category", http.StatusBadRequest)
		return
	}
	utils.MarshalAndWriteResponse(w, category)
}

func (a *API) PostSub(w http.ResponseWriter, r *http.Request) {
	var request PostCategoriesRequestBody
	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}
	category, err := a.categoryManagment.CreateCategory(r.Context(), services.CreateCategoryRequest{
		CompanyID:   request.CompanyID,
		CategoryID:  request.CategoryID,
		Description: request.Description,
		Type:        request.Type})
	if err != nil {
		log.Printf("Error Creating a Category: %v", err)
		http.Error(w, "Error Creating a Category", http.StatusBadRequest)
		return
	}
	utils.MarshalAndWriteResponse(w, category)
}

func (a *API) PutCategories(w http.ResponseWriter, r *http.Request) {
	var request PutCategoriesRequestBody
	err := utils.UnmarshalRequest(r, &request)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}
	category, err := a.categoryManagment.UpdateCategory(r.Context(), request.CategoryID, services.UpdateCategoryRequest{
		CategoryID:  request.CategoryID,
		Description: request.Description})
	if err != nil {
		log.Printf("Error Updating a Category: %v", err)
		http.Error(w, "Error Updating a Category", http.StatusBadRequest)
		return
	}
	utils.MarshalAndWriteResponse(w, category)
}

func (a *API) GetCategories(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["id"]
	categories, err := a.categoryManagment.GetCategory(r.Context(), companyId)
	if err != nil {
		log.Printf("Error Getting Categories: %v", err)
		http.Error(w, "Error Getting Categories", http.StatusBadRequest)
		return
	}
	utils.MarshalAndWriteResponse(w, GetCategoriesResponseBody{
		TotalCompanies: len(categories),
		Description:    categories,
	})
}

func (a *API) GetSubCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := a.categoryManagment.GetSub(r.Context(), r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Error Getting Sub Categories: %v", err)
		http.Error(w, "Error Getting Sub Categories", http.StatusBadRequest)
		return
	}
	utils.MarshalAndWriteResponse(w, GetCategoriesResponseBody{
		TotalCompanies: len(categories),
		Description:    categories,
	})
}

func (a *API) DeleteCategories(w http.ResponseWriter, r *http.Request) {
	category, err := a.categoryManagment.DeleteCategory(r.Context(), r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Error Deleting a Category: %v", err)
		http.Error(w, "Error Deleting a Category", http.StatusBadRequest)
		return
	}
	utils.MarshalAndWriteResponse(w, category)
}
