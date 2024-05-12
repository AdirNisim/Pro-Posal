package models

import "time"

type Company struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	ContactID  string    `json:"contact_id"`
	Address    string    `json:"address"`
	LogoBase64 string    `json:"logo_base64"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeleteAt   time.Time `json:"deleted_at"`
}
