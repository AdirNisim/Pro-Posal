package models

import "time"

type Role string

const (
	AdminRole                 Role = "admin"
	CompanyAdminRole          Role = "company_admin"
	CompanyContributorRole    Role = "company_contributor"
	CompanyProjectManagerRole Role = "company_project_manager"
	ProspectRole              Role = "prospect"
)

type Permission struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CompanyID  string    `json:"company_id"`
	Role       Role      `json:"role"`
	ContractID string    `json:"contract_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserPermission struct {
	User       User
	Permission Permission
}

/*
"admin" - Admin is the perspon that perchase the subscription and own the account
		  Admin can do everything within the account realm.

"company_admin" - CompanyAdmin is the person that is in charge of the company (for Example CEO)
				  he can do everything within the company realm
				  Exceptions : * Can not Delete the company
				  			   * Can not Create new Users

"company_contributor" - CompanyContributor is the person that is authorized make changes and create Templates and Offers

"company_project_manager" - CompanyProjectManager is the sales person in the field that is authorized to create new offers from the pre-defined templates.
							Need to think on: CompanyProjectManager can create/Invite a new prospact user.

"prospect" - Prospect is the person that will receive the offer and can accept or reject the offer.
*/
