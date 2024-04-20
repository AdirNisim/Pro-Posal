package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pro-posal/webserver/dao"
	"github.com/pro-posal/webserver/internal/database"
	"github.com/pro-posal/webserver/internal/utils"
	"github.com/pro-posal/webserver/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Password  string
	InvitedBy *string
}

type UserManagementService interface {
	CreateUser(context.Context, CreateUserRequest) (*models.User, error)
	ListUsers(context.Context) ([]*models.User, error)
}

type userManagementServiceImpl struct {
	db *database.DBConnector
}

func NewUserManagementService(db *database.DBConnector) UserManagementService {
	return &userManagementServiceImpl{
		db: db,
	}
}

func (s *userManagementServiceImpl) CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed hashing password: %w", err)
	}

	userDao := dao.User{
		ID:           uuid.NewString(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Email:        req.Email,
		EmailHash:    utils.HashEmail(req.Email),
		PasswordHash: hashedPassword,
		InvitedBy:    null.StringFromPtr(req.InvitedBy),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// TODO: Check if already exists to ensure non-500 errors for idempotency

	err = userDao.Insert(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("failed inserting user to database: %w", err)
	}

	return userDaoToUserModel(userDao), nil
}

func (s *userManagementServiceImpl) ListUsers(ctx context.Context) ([]*models.User, error) {
	userDaos, err := dao.Users().All(ctx, s.db.Conn)
	if err != nil {
		return nil, fmt.Errorf("failed fetching users from database: %w", err)
	}

	var users []*models.User
	for _, userDao := range userDaos {
		users = append(users, userDaoToUserModel(*userDao))
	}

	return users, nil
}

func userDaoToUserModel(userDao dao.User) *models.User {
	return &models.User{
		ID:        userDao.ID,
		FirstName: userDao.FirstName,
		LastName:  userDao.LastName,
		Phone:     userDao.Phone,
		Email:     userDao.Email,
		CreatedAt: userDao.CreatedAt,
		UpdatedAt: userDao.UpdatedAt,
	}
}
