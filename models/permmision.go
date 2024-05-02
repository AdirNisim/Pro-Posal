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
