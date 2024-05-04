package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pro-posal/webserver/config"
	"github.com/pro-posal/webserver/dao"
	"github.com/pro-posal/webserver/internal/database"
	"github.com/pro-posal/webserver/internal/utils"
	"github.com/pro-posal/webserver/models"
)

type AuthService interface {
	CreateAuthToken(ctx context.Context, email string, password string) (*models.AuthToken, error)
	ValidateAuthToken(context.Context, string) (*models.Session, error)
}

type authServiceImpl struct {
	db *database.DBConnector
}

func NewAuthService(db *database.DBConnector) AuthService {
	return &authServiceImpl{
		db: db,
	}
}

func (s *authServiceImpl) CreateAuthToken(ctx context.Context, email string, password string) (*models.AuthToken, error) {
	userDao, err := dao.Users(
		dao.UserWhere.EmailHash.EQ(utils.HashEmail(email)),
		// dao.UserWhere.DeletedAt.IsNull(), // TODO - Example when you'll introduce deleted_at
	).One(ctx, s.db.Conn)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("failed fetching user from database: %w", err)
	}

	if !utils.ComparePasswords(userDao.PasswordHash, password) {
		return nil, errors.New("invalid email or password")
	}

	userID, err := uuid.Parse(userDao.ID)
	if err != nil {
		return nil, fmt.Errorf("failed parsing user ID %v: %w", userDao.ID, err) // Should never happen
	}

	return s.createTokenFromSession(ctx, &models.Session{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(config.AppConfig.Auth.ExpirationTimeMinutes) * time.Minute),
	})
}

func (s *authServiceImpl) ValidateAuthToken(ctx context.Context, token string) (*models.Session, error) {
	at, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.Auth.JWTSigningSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed parsing token: %w", err)
	}

	if !at.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := at.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims within token")
	}

	// TODO - For options #2 - you'll just get the sessionID and fetch the session object from the database
	// createdAt, err := time.Parse("2006-01-02T15:04:05.000Z", claims["session"].(map[string]any)["created_at"].(string))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed parsing created_at from claims: %w", err)
	// }
	// expiresAt, err := time.Parse("2006-01-02T15:04:05.000Z", claims["session"].(map[string]any)["expires_at"].(string))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed parsing expires_at from claims: %w", err)
	// }
	session := &models.Session{
		ID:        uuid.MustParse(claims["session"].(map[string]any)["id"].(string)),
		UserID:    uuid.MustParse(claims["session"].(map[string]any)["user_id"].(string)),
		// CreatedAt: createdAt,
		// ExpiresAt: expiresAt,
	}

	// Sanity Check
	if session.UserID.String() != claims["sub"].(string) {
		return nil, errors.New("malformed session in token")
	}

	return session, nil
}

func (s *authServiceImpl) createTokenFromSession(_ context.Context, session *models.Session) (*models.AuthToken, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": session.UserID,
		// "session", session.ID, // TODO - This is for option #2
		"session": session,
		"exp":     session.ExpiresAt.Unix(),
	})

	bearerToken, err := at.SignedString([]byte(config.AppConfig.Auth.JWTSigningSecret))
	if err != nil {
		return nil, fmt.Errorf("failed signing token: %w", err)
	}

	return &models.AuthToken{
		BearerToken: bearerToken,
		ExpiresAt:   session.ExpiresAt,
	}, nil
}
