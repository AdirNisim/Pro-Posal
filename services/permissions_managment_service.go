package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pro-posal/webserver/dao"
	"github.com/pro-posal/webserver/internal/database"
	"github.com/pro-posal/webserver/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type CreatePermissionRequest struct {
	UserID     string
	CompanyID  string
	Role       string
	ContractID string
}

type UpdatePermissionRequest struct {
	UserID    string
	CompanyID string
	Role      string
}

type PermissionManagementService interface {
	CreatePermission(context.Context, CreatePermissionRequest) (*models.Permission, error)
	UpdatePermission(context.Context, UpdatePermissionRequest) (*models.Permission, error)
}

type PermissionManagementServiceImpl struct {
	db *database.DBConnector
}

func NewPermissionManagementService(db *database.DBConnector) PermissionManagementService {
	return &PermissionManagementServiceImpl{
		db: db,
	}
}

func (s *PermissionManagementServiceImpl) CreatePermission(ctx context.Context, req CreatePermissionRequest) (*models.Permission, error) {
	permissionDao := dao.Permission{
		ID:         uuid.NewString(),
		UserID:     req.UserID,
		CompanyID:  req.CompanyID,
		Role:       req.Role,
		ContractID: req.ContractID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := permissionDao.Insert(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("failed inserting permission to database: %w", err)
	}

	return permissionsDaoToPermissionModel(permissionDao), nil
}

func (s *PermissionManagementServiceImpl) UpdatePermission(ctx context.Context, req UpdatePermissionRequest) (*models.Permission, error) {
	// Find the existing permission based on UserID and CompanyID
	permission, err := dao.Permissions(
		qm.Where("user_id = ? AND company_id = ?", req.UserID, req.CompanyID),
	).One(ctx, s.db.Conn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no permission found with UserID %s and CompanyID %s", req.UserID, req.CompanyID)
		}
		return nil, fmt.Errorf("error fetching permission: %w", err)
	}

	permission.Role = req.Role
	permission.UpdatedAt = time.Now()

	_, err = permission.Update(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("error updating permission: %w", err)
	}

	updatedPermission := permissionsDaoToPermissionModel(*permission)
	return updatedPermission, nil
}

func permissionsDaoToPermissionModel(permissionDao dao.Permission) *models.Permission {
	return &models.Permission{
		ID:         permissionDao.ID,
		UserID:     permissionDao.UserID,
		CompanyID:  permissionDao.CompanyID,
		Role:       permissionDao.Role,
		ContractID: permissionDao.ContractID,
		CreatedAt:  permissionDao.CreatedAt,
		UpdatedAt:  permissionDao.UpdatedAt,
	}
}
