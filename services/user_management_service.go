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

type ChangeUserPasswordRequest struct {
	Uuid        string
	NewPassword string
}

type UserManagementService interface {
	CreateUser(context.Context, CreateUserRequest) (*models.User, error)
	ListUsers(context.Context) ([]*models.User, error)
	GetUserByID(context.Context, string) (*models.User, error)
	UpdateUserPassword(context.Context, ChangeUserPasswordRequest) (*models.User, error)
	DeleteUser(context.Context, string) (*models.User, error)
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

	// Check if already exists to ensure non-500 errors for idempotency
	existingUser, err := dao.Users(dao.UserWhere.Email.EQ(req.Email)).One(ctx, s.db.Conn)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed checking for existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

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

func (s *userManagementServiceImpl) GetUserByID(ctx context.Context, userId string) (*models.User, error) {
	userObj, err := dao.FindUser(ctx, s.db.Conn, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return userDaoToUserModel(*userObj), nil
}

func (s *userManagementServiceImpl) UpdateUserPassword(ctx context.Context, req ChangeUserPasswordRequest) (*models.User, error) {
	userDao, err := dao.FindUser(ctx, s.db.Conn, req.Uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("failed hashing new password: %w", err)
	}

	userDao.PasswordHash = hashedPassword
	userDao.UpdatedAt = time.Now()

	_, err = userDao.Update(ctx, s.db.Conn, boil.Whitelist("password_hash", "updated_at"))
	if err != nil {
		return nil, fmt.Errorf("failed to update user password in database: %w", err)
	}

	updatedUser := userDaoToUserModel(*userDao)
	return updatedUser, nil
}

func (s *userManagementServiceImpl) DeleteUser(ctx context.Context, userId string) (*models.User, error) {
	userDao, err := dao.FindUser(ctx, s.db.Conn, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	deletedAt := null.TimeFrom(time.Now())
	userDao.DeletedAt = deletedAt

	_, err = userDao.Update(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("error updating company: %w", err)
	}

	return userDaoToUserModel(*userDao), nil

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
