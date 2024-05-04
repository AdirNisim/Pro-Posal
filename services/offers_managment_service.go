package services

import "github.com/pro-posal/webserver/internal/database"

type OfferManagementService interface {
}

type OfferManagementServiceImpl struct {
	db *database.DBConnector
}

func NewOfferManagementService(db *database.DBConnector) OfferManagementService {
	return &OfferManagementServiceImpl{
		db: db,
	}
}
