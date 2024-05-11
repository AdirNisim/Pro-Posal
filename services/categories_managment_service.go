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

type CreateCategoryRequest struct {
	CompanyID   string
	CategoryID  string
	Description string
	Type        string
}

type UpdateCategoryRequest struct {
	CategoryID  string
	Description string
}

type CategoryManagementService interface {
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*models.Category, error)
	DeleteCategory(ctx context.Context, id string) (*models.Category, error)
	UpdateCategory(ctx context.Context, id string, req UpdateCategoryRequest) (*models.Category, error)
	GetCategory(ctx context.Context, companyID string) ([]*models.Category, error)
	GetSubCategory(ctx context.Context, id string) ([]*models.Category, error)
	GetDescription(ctx context.Context, id string) ([]*models.Category, error)
}

type CategoryManagementServiceImpl struct {
	db *database.DBConnector
}

func NewCategoryManagementService(db *database.DBConnector) CategoryManagementService {
	return &CategoryManagementServiceImpl{
		db: db,
	}
}

func (s *CategoryManagementServiceImpl) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*models.Category, error) {
	existingCategory, err := dao.Categories(qm.Where("description= ? AND company_id = ?", req.Description, req.CompanyID)).One(ctx, s.db.Conn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No existing category found; continue to create a new one.
		} else {
			return nil, fmt.Errorf("error checking if category exists: %w", err)
		}
	} else if existingCategory != nil {
		return nil, fmt.Errorf("category already exists")
	}

	categoryDao := dao.Category{
		ID:          uuid.NewString(),
		CategoryID:  null.String{String: req.CategoryID, Valid: true},
		Description: req.Description,
		CompanyID:   req.CompanyID,
		Type:        req.Type,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = categoryDao.Insert(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("failed to insert category into database: %w", err)
	}
	return categoryDaoToCategoryModel(categoryDao), nil
}

func (s *CategoryManagementServiceImpl) UpdateCategory(ctx context.Context, id string, req UpdateCategoryRequest) (*models.Category, error) {
	categoryDao, err := dao.Categories(qm.Where("id = ?", id)).One(ctx, s.db.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed to get category from database: %w", err)
	}

	categoryDao.Description = req.Description
	categoryDao.UpdatedAt = time.Now()

	_, err = categoryDao.Update(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("failed to update category in database: %w", err)
	}

	return categoryDaoToCategoryModel(*categoryDao), nil
}

func (s *CategoryManagementServiceImpl) DeleteCategory(ctx context.Context, id string) (*models.Category, error) {
	categoryDao, err := dao.Categories(qm.Where("id = ?", id)).One(ctx, s.db.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed to get category from database: %w", err)
	}
	if categoryDao == nil {
		return nil, fmt.Errorf("category not found")
	}
	if categoryDao.Type == "description" {
		return s.deleteCategoryHelper(ctx, categoryDao)
	} else {
		subDao, err := dao.Categories(qm.Where("category_id = ?", categoryDao.ID)).All(ctx, s.db.Conn)
		if err != nil {
			return nil, fmt.Errorf("failed to get sub categories from database: %w", err)
		}
		for _, sub := range subDao {
			s.DeleteCategory(ctx, sub.ID)
		}
	}
	return s.deleteCategoryHelper(ctx, categoryDao)
}

func (s *CategoryManagementServiceImpl) deleteCategoryHelper(ctx context.Context, categoryDao *dao.Category) (*models.Category, error) {
	_, err := categoryDao.Delete(ctx, s.db.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed to delete category from database: %w", err)
	}
	return categoryDaoToCategoryModel(*categoryDao), nil
}

func (s *CategoryManagementServiceImpl) GetCategory(ctx context.Context, companyID string) ([]*models.Category, error) {
	categoryDao, err := dao.Categories(qm.Where("company_id = ? AND type = ?", companyID, "category")).All(ctx, s.db.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories from database: %w", err)
	}
	var categories []*models.Category
	for _, category := range categoryDao {
		categories = append(categories, categoryDaoToCategoryModel(*category))
	}
	return categories, nil
}

func (s *CategoryManagementServiceImpl) GetSubCategory(ctx context.Context, id string) ([]*models.Category, error) {
	categoryDao, err := dao.Categories(qm.Where("category_id = ?", id)).All(ctx, s.db.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed to get sub categories from database: %w", err)
	}
	var subcategories []*models.Category
	for _, category := range categoryDao {
		subcategories = append(subcategories, categoryDaoToCategoryModel(*category))
	}
	return subcategories, nil
}

func (s *CategoryManagementServiceImpl) GetDescription(ctx context.Context, id string) ([]*models.Category, error) {
	categoryDao, err := dao.Categories(qm.Where("company_id = ?", id)).All(ctx, s.db.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed to get descriptions from database: %w", err)
	}
	var descriptions []*models.Category
	for _, subcategory := range categoryDao {
		descriptions = append(descriptions, categoryDaoToCategoryModel(*subcategory))
	}
	return descriptions, nil

}

func categoryDaoToCategoryModel(categoryDao dao.Category) *models.Category {
	return &models.Category{
		ID:          categoryDao.ID,
		CompanyID:   categoryDao.CompanyID,
		CategoryID:  categoryDao.CategoryID.String,
		Description: categoryDao.Description,
		CreatedAt:   categoryDao.CreatedAt,
		UpdatedAt:   categoryDao.UpdatedAt,
	}
}
