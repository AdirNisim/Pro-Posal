package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pro-posal/webserver/dao"
	"github.com/pro-posal/webserver/internal/database"
	"github.com/pro-posal/webserver/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type CreateContractTemplateRequest struct {
	Name      string
	CompanyID string
	Template  string
}

type UpdateContractsTemplatesRequest struct {
	Name     string
	Template string
}

type ContractTemplateManagementService interface {
	PostContractsTemplates(context.Context, CreateContractTemplateRequest) (*models.ContractTemplate, error)
	GetContractsTemplates(context.Context, string) (*models.ContractTemplate, error)
	GetContractsTemplate(context.Context, string) ([]*models.ContractTemplate, error)
	UpdateContractsTemplate(context.Context, string, UpdateContractsTemplatesRequest) (*models.ContractTemplate, error)
	DeleteContractsTemplates(context.Context, string) (*models.ContractTemplate, error)
}

type ContractTemplateManagementServiceImpl struct {
	db *database.DBConnector
}

func NewContractTemplateManagementService(db *database.DBConnector) ContractTemplateManagementService {
	return &ContractTemplateManagementServiceImpl{
		db: db,
	}
}

func (s *ContractTemplateManagementServiceImpl) PostContractsTemplates(ctx context.Context, req CreateContractTemplateRequest) (*models.ContractTemplate, error) {
	existingContract, err := dao.ContractTemplates(
		qm.Where("name = ? AND company_id = ?", req.Name, req.CompanyID),
	).One(ctx, s.db.Conn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No existing contract template found; continue to create a new one.
		} else {
			return nil, fmt.Errorf("error checking if contract template exists: %w", err)
		}
	} else if existingContract != nil {
		return nil, fmt.Errorf("contract template already exists")
	}

	contractDao := dao.ContractTemplate{
		ID:        uuid.NewString(),
		Name:      req.Name,
		CompanyID: req.CompanyID,
		Template:  req.Template,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = contractDao.Insert(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("failed to insert contract template into database: %w", err)
	}

	return contractDaoToContractModel(contractDao), nil
}

func contractDaoToContractModel(contractDao dao.ContractTemplate) *models.ContractTemplate {
	return &models.ContractTemplate{
		ID:        contractDao.ID,
		Name:      contractDao.Name,
		CompanyID: contractDao.CompanyID,
		Template:  contractDao.Template,
		CreatedAt: contractDao.CreatedAt,
		UpdatedAt: contractDao.UpdatedAt,
	}
}
