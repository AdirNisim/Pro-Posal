package api

import "net/http"

type PostPremmisionRequestBody struct {
	UserID     string `json:"user_id"`
	CompanyID  string `json:"company_id"`
	Role       string `json:"role"`
	ContractID string `json:"contract_id"`
}

func (a *API) PostPermmision(w http.ResponseWriter, r *http.Request) {

}
