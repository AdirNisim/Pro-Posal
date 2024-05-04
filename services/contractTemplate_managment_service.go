package services

import "github.com/pro-posal/webserver/internal/database"

type ContractTemplateManagementService interface {
}

type ContractTemplateManagementServiceImpl struct {
	db *database.DBConnector
}

func NewContractTemplateManagementService(db *database.DBConnector) ContractTemplateManagementService {
	return &ContractTemplateManagementServiceImpl{
		db: db,
	}
}
