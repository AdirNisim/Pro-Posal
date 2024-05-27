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

type CreatePermissionRequest struct {
	UserID     string
	CompanyID  string
	Role       string
	ContractID string
}

type UpdatePermissionRequest struct {
	Id         string
	Role       string
	ContractID string
}

type PermissionManagementService interface {
	CreatePermission(context.Context, CreatePermissionRequest) (*models.Permission, error)
	UpdatePermission(context.Context, UpdatePermissionRequest) (*models.Permission, error)
	GetPermissions(context.Context, string) ([]*models.UserPermission, error)
	DeletePermission(context.Context, string) (*models.Permission, error)
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

	return permissionDaoToPermissionModel(&permissionDao), nil
}

func (s *PermissionManagementServiceImpl) UpdatePermission(ctx context.Context, req UpdatePermissionRequest) (*models.Permission, error) {
	permission, err := dao.Permissions(
		qm.Where("id = ?", req.Id),
	).One(ctx, s.db.Conn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no permission found with ID %s", req.Id)
		}
		return nil, fmt.Errorf("error fetching permission: %w", err)
	}

	permission.Role = req.Role
	permission.ContractID = req.ContractID
	permission.UpdatedAt = time.Now()

	_, err = permission.Update(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("error updating permission: %w", err)
	}

	updatedPermission := permissionDaoToPermissionModel(permission)
	return updatedPermission, nil
}

func (s *PermissionManagementServiceImpl) GetPermissions(ctx context.Context, companyId string) ([]*models.UserPermission, error) {
	permissionsDao, err := dao.Permissions(
		qm.Select("permissions.*, users.*"),
		qm.InnerJoin("users on permissions.user_id = users.id"),
		qm.Where("permissions.company_id = ?", companyId),
	).All(ctx, s.db.Conn)
	if err != nil {
		return nil, fmt.Errorf("error fetching permissions: %w", err)
	}

	userPermissions := make([]*models.UserPermission, len(permissionsDao))
	for i, permissionDao := range permissionsDao {
		userPermissions[i] = &models.UserPermission{
			User: models.User{
				ID:        permissionDao.R.User.ID,
				FirstName: permissionDao.R.User.FirstName,
				LastName:  permissionDao.R.User.LastName,
				Phone:     permissionDao.R.User.Phone,
				Email:     permissionDao.R.User.Email,
				CreatedAt: permissionDao.R.User.CreatedAt,
				UpdatedAt: permissionDao.R.User.UpdatedAt,
			},
			Permission: *permissionDaoToPermissionModel(permissionDao),
		}
	}

	return userPermissions, nil
}

func (s *PermissionManagementServiceImpl) DeletePermission(ctx context.Context, id string) (*models.Permission, error) {

	permissionDao, err := dao.FindPermission(ctx, s.db.Conn, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no permission found with ID %s", id)
		}
		return nil, fmt.Errorf("error retrieving permission: %w", err)
	}
	permission := permissionDaoToPermissionModel(permissionDao)

	_, err = permissionDao.Delete(ctx, s.db.Conn)
	if err != nil {
		return nil, fmt.Errorf("error deleting permission: %w", err)
	}

	return permission, nil
}

func permissionDaoToPermissionModel(permissionDao *dao.Permission) *models.Permission {
	return &models.Permission{
		ID:         permissionDao.ID,
		UserID:     permissionDao.UserID,
		CompanyID:  permissionDao.CompanyID,
		Role:       models.Role(permissionDao.Role),
		ContractID: permissionDao.ContractID,
		CreatedAt:  permissionDao.CreatedAt,
		UpdatedAt:  permissionDao.UpdatedAt,
	}
}
