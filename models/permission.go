package models

import "time"

type Permission struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CompanyID  string    `json:"company_id"`
	Role       string    `json:"role"`
	ContractID string    `json:"contract_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// When you will implement your authorization middleware
// Then for each request you will authorize it based on your structure
// For example, if you are role based
// if method == "POST" && path == "/companies" {
//	return user.Role == "admin"
// }

// Or if you are resource/action based
// if method == "POST" && path == "/companies" {
//	return user.permissions.CanManageCompanies()
// }

// Also, you can store a map of allowed operations in the user struct, per API
// user.Permissions := map[string][]string{
// "POST", "/companies",
// "GET", "/offers",
// }
