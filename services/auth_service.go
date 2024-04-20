package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pro-posal/webserver/dao"
	"github.com/pro-posal/webserver/internal/database"
	"github.com/pro-posal/webserver/internal/utils"
)

type AuthService interface {
	CreateAuthToken(ctx context.Context, email string, password string) (string, time.Time, error)
	ValidateAuthToken(context.Context, string) error
}

type authServiceImpl struct {
	db *database.DBConnector
}

func NewAuthService(db *database.DBConnector) AuthService {
	return &authServiceImpl{
		db: db,
	}
}

func (s *authServiceImpl) CreateAuthToken(ctx context.Context, email string, password string) (string, time.Time, error) {
	userDao, err := dao.Users(
		dao.UserWhere.EmailHash.EQ(utils.HashEmail(email)),
	).One(ctx, s.db.Conn)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", time.Now(), errors.New("invalid email or password")
		}
		return "", time.Now(), fmt.Errorf("failed fetching user from database: %w", err)
	}

	if !utils.ComparePasswords(userDao.PasswordHash, password) {
		return "", time.Now(), errors.New("invalid email or password")
	}

	// TODO: Generate expiring token

	return "you-got-access", time.Now().Add(5 * time.Minute), nil
}

func (s *authServiceImpl) ValidateAuthToken(ctx context.Context, token string) error {
	// TODO: Real implementation with JWT
	if token != "you-got-access" {
		return errors.New("invalid token")
	}
	return nil
}
