package services

import "github.com/pro-posal/webserver/internal/database"

type CategoryManagementService interface {
}

type CategoryManagementServiceImpl struct {
	db *database.DBConnector
}

func NewCategoryManagementService(db *database.DBConnector) CategoryManagementService {
	return &CategoryManagementServiceImpl{
		db: db,
	}
}
