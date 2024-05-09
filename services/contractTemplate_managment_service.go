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
	"github.com/volatiletech/null/v8"
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
	PostContractsTemplate(context.Context, CreateContractTemplateRequest) (*models.ContractTemplate, error)
	GetContractsTemplate(context.Context, string) (*models.ContractTemplate, error)
	GetContractsTemplates(context.Context, string) ([]*models.ContractTemplate, error)
	UpdateContractsTemplate(context.Context, string, UpdateContractsTemplatesRequest) (*models.ContractTemplate, error)
	DeleteContractsTemplate(context.Context, string) (*models.ContractTemplate, error)
}

type ContractTemplateManagementServiceImpl struct {
	db *database.DBConnector
}

func NewContractTemplateManagementService(db *database.DBConnector) ContractTemplateManagementService {
	return &ContractTemplateManagementServiceImpl{
		db: db,
	}
}

func (s *ContractTemplateManagementServiceImpl) PostContractsTemplate(ctx context.Context, req CreateContractTemplateRequest) (*models.ContractTemplate, error) {
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

func (s *ContractTemplateManagementServiceImpl) DeleteContractsTemplate(ctx context.Context, id string) (*models.ContractTemplate, error) {
	contractTemplateDoa, err := dao.FindContractTemplate(ctx, s.db.Conn, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no contract template found with ID %s", id)
		}
		return nil, fmt.Errorf("error retrieving contract template: %w", err)
	}
	deletedAt := null.TimeFrom(time.Now())
	contractTemplateDoa.DeletedAt = deletedAt

	contract := contractDaoToContractModel(*contractTemplateDoa)

	_, err = contractTemplateDoa.Update(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("error deleteing contract template: %w", err)
	}

	return contract, nil
}

func (s *ContractTemplateManagementServiceImpl) UpdateContractsTemplate(ctx context.Context, id string, req UpdateContractsTemplatesRequest) (*models.ContractTemplate, error) {
	contractTemplateDoa, err := dao.FindContractTemplate(ctx, s.db.Conn, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no contract template found with ID %s", id)
		}
		return nil, fmt.Errorf("error retrieving contract template: %w", err)
	}

	contractTemplateDoa.Name = req.Name
	contractTemplateDoa.Template = req.Template
	contractTemplateDoa.UpdatedAt = time.Now()

	_, err = contractTemplateDoa.Update(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("error updating contract template: %w", err)
	}

	return contractDaoToContractModel(*contractTemplateDoa), nil
}

func (s *ContractTemplateManagementServiceImpl) GetContractsTemplate(ctx context.Context, id string) (*models.ContractTemplate, error) {
	contractTemplateDoa, err := dao.FindContractTemplate(ctx, s.db.Conn, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no contract template found with ID %s", id)
		}
		return nil, fmt.Errorf("error retrieving contract template: %w", err)
	}

	return contractDaoToContractModel(*contractTemplateDoa), nil
}

func (s *ContractTemplateManagementServiceImpl) GetContractsTemplates(ctx context.Context, companyID string) ([]*models.ContractTemplate, error) {
	contractTemplates, err := dao.ContractTemplates(
		qm.Where("company_id = ?", companyID),
	).All(ctx, s.db.Conn)
	if err != nil {
		return nil, fmt.Errorf("error retrieving contract templates: %w", err)
	}

	var contractTemplateModels []*models.ContractTemplate
	for _, contractTemplate := range contractTemplates {
		contractTemplateModels = append(contractTemplateModels, contractDaoToContractModel(*contractTemplate))
	}

	return contractTemplateModels, nil
}
