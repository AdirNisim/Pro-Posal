package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pro-posal/webserver/dao"
	"github.com/pro-posal/webserver/internal/database"
	"github.com/pro-posal/webserver/internal/utils"
	"github.com/pro-posal/webserver/models"
	"github.com/pro-posal/webserver/services"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func AuthorizationMiddleware(authService services.AuthService, next http.Handler, db *database.DBConnector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := utils.GetSessionFromContext(r.Context())
		if session == nil {
			log.Println("No session found in the context, skipping authorization")
			next.ServeHTTP(w, r)
			return
		}

		permissionsDao, err := dao.Permissions(qm.Where("user_id = ?", session.UserID)).All(r.Context(), db.Conn)
		if err != nil {
			log.Print("Error fetching permissions from the database: ", err)
			next.ServeHTTP(w, r)
			return
		}

		permissions := make([]*models.Permission, len(permissionsDao))
		for i, permissionDao := range permissionsDao {
			permissions[i] = &models.Permission{
				ID:         permissionDao.ID,
				UserID:     permissionDao.UserID,
				CompanyID:  permissionDao.CompanyID,
				Role:       models.Role(permissionDao.Role),
				ContractID: permissionDao.ContractID,
				CreatedAt:  permissionDao.CreatedAt,
				UpdatedAt:  permissionDao.UpdatedAt,
			}
		}

		if hasAdminRole(permissions) {
			log.Printf("Authorized admin %v to call %v", session.UserID, r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}

		pathParts := strings.Split(r.URL.Path, "/")
		resourceIndex := 0

		for resourceIndex < len(pathParts) {
			hasMoreParts := resourceIndex+2 < len(pathParts)

			switch pathParts[resourceIndex] {
			case "companies":
				companyID := mux.Vars(r)["companyId"]
				err = authService.AuthorizeCompany(r.Method, session.UserID, permissions, hasMoreParts, companyID)

			// case "users":
			// 	userID := mux.Vars(r)["userId"]
			// 	err = authService.AuthorizeUser(r.Method, session.UserID, userID, hasMoreParts)

			case "contracts":
				companyID := mux.Vars(r)["companyId"]
				contractID := mux.Vars(r)["contractId"]
				err = authService.AuthorizeContract(r.Method, session.UserID, permissions, hasMoreParts, companyID, contractID)
			}

			if err != nil {
				if _, ok := err.(*services.UnauthorizedError); ok {
					log.Printf("Unauthorized user %v tried to call %v", session.UserID, r.URL.Path)
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				log.Printf("Error authorizing user %v: %v", session.UserID, err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			resourceIndex += 2
		}

		log.Printf("Authorized user %v to call %v", session.UserID, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func hasAdminRole(userPermissions []*models.Permission) bool {
	for _, permission := range userPermissions {
		if permission.Role == models.AdminRole {
			return true
		}
	}
	return false
}
