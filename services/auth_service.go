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
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UnauthorizedError struct{}

func (u *UnauthorizedError) Error() string {
	return "Unauthorized"
}

type AuthService interface {
	CreateAuthToken(ctx context.Context, email string, password string) (*models.AuthToken, error)
	ValidateAuthToken(context.Context, string) (*models.Session, error)

	AuthorizeCompany(method string, callerUserID uuid.UUID, permissions []*models.Permission, hasMoreParts bool, requestedCompanyId string) error
	AuthorizeContract(method string, callerUserID uuid.UUID, permissions []*models.Permission, hasMoreParts bool, requestedCompanyId string, requestedContractId string) error
	// AuthorizeUser(method string, callerUserID uuid.UUID, requestedUserId string) error
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
		dao.UserWhere.DeletedAt.IsNull(),
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

	session := models.Session{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(config.AppConfig.Auth.ExpirationTimeMinutes) * time.Minute),
	}

	sessionDao := dao.Session{
		ID:        session.ID.String(),
		UserID:    session.UserID.String(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(config.AppConfig.Auth.ExpirationTimeMinutes) * time.Minute),
	}

	err = sessionDao.Insert(ctx, s.db.Conn, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("failed to insert session into database: %w", err)
	}

	return s.createTokenFromSession(ctx, &session)
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
		ID:     uuid.MustParse(claims["session"].(map[string]any)["id"].(string)),
		UserID: uuid.MustParse(claims["session"].(map[string]any)["user_id"].(string)),
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

func (s *authServiceImpl) AuthorizeCompany(method string, callerUserID uuid.UUID, permissions []*models.Permission, hasMoreParts bool, requestedCompanyId string) error {
	// Example for DB table entries
	// | user_id | role | company_id | contract_id
	// | blah1   | admin | - | -
	// | blah2   | company_admin | company1 | -
	// | blah2   | company_project_manager | company5 | -
	// | blah3   | company_contributor | company1 | -
	// | blah4   | prospect | company1 | contract1
	// | blah4   | prospect | company1 | contract2

	if method == "POST" {
		// POST /companies/{companyId}/contracts
		// POST /companies/{companyId}/categories
		if hasMoreParts {
			return callerHasRolesForCompanyById(requestedCompanyId, permissions, models.CompanyAdminRole, models.CompanyProjectManagerRole, models.CompanyContributorRole)
		}

		// POST /companies
		return &UnauthorizedError{}
	}

	if method == "PUT" {
		// PUT /companies/{companyId}/contracts/{contractId}
		if hasMoreParts {
			return callerHasRolesForCompanyById(requestedCompanyId, permissions, models.CompanyAdminRole, models.CompanyProjectManagerRole, models.CompanyContributorRole)
		}

		// PUT /companies/{companyId}
		return callerHasRolesForCompanyById(requestedCompanyId, permissions, models.CompanyAdminRole)
	}

	if method == "GET" {
		// GET /companies/{companyId}/...
		return callerHasRolesForCompanyById(requestedCompanyId, permissions, models.CompanyAdminRole, models.CompanyProjectManagerRole, models.CompanyContributorRole)
	}

	// Unsupported method
	return &UnauthorizedError{}
}

func (s *authServiceImpl) AuthorizeContract(method string, callerUserID uuid.UUID, permissions []*models.Permission, hasMoreParts bool, requestedCompanyId string, requestedContractId string) error {
	if method == "POST" {
		// POST /companies/{companyId}/contracts
		return callerHasRolesForCompanyById(requestedCompanyId, permissions, models.CompanyProjectManagerRole)
	}

	if method == "GET" {
		// TODO: Implement this
	}

	// Unsupported method
	return &UnauthorizedError{}
}

func callerHasRolesForCompanyById(requestedCompanyId string, permissions []*models.Permission, allowedRoles ...models.Role) error {
	for _, permission := range permissions {
		if permission.CompanyID == requestedCompanyId {
			for _, role := range allowedRoles {
				if permission.Role == role {
					return nil
				}
			}
		}
	}
	return &UnauthorizedError{}
}

func callerHasRolesForContractById(requestedCompanyId string, requestedContractId string, permissions []*models.Permission, allowedRoles ...models.Role) error {
	for _, permission := range permissions {
		if permission.CompanyID == requestedCompanyId && permission.ContractID == requestedContractId {
			for _, role := range allowedRoles {
				if permission.Role == role {
					return nil
				}
			}
		}
	}
	return &UnauthorizedError{}
}
