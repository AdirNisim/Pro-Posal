package models

import "time"

type ContractTemplate struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CompanyID string    `json:"company_id"`
	Template  string    `json:"template"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
