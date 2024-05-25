package models

import "time"

type CategoryType string

const (
	CategoryTypeCategory    CategoryType = "category"
	CategoryTypeSubCategory CategoryType = "sub_category"
	CategoryTypeDescription CategoryType = "description"
)

type Category struct {
	ID          string       `json:"id"`
	CompanyID   string       `json:"company_id"`
	CategoryID  string       `json:"category_id"`
	Description string       `json:"description"`
	Type        CategoryType `json:"type"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeleteAt    time.Time    `json:"deleted_at"`
}
