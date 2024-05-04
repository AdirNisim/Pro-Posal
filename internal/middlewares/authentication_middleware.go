package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/pro-posal/webserver/services"
)

var AuthBypassRoutes = map[string]string{
	"/users/login": "POST",
}

func AuthenticationMiddleware(authService services.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Bypass authentication if needed
		if _, ok := AuthBypassRoutes[r.URL.Path]; ok {
			if r.Method == AuthBypassRoutes[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Perform authentication
		authToken := extractBearerToken(w, r)
		if authToken == "" {
			return
		}

		session, err := authService.ValidateAuthToken(r.Context(), authToken)
		if err != nil {
			log.Printf("Failed validating auth token: %v", err)
			http.Error(w, "Invalid Bearer Token", http.StatusForbidden)
			return
		}

		// Add the session to the context, so it can be used within different handlers
		ctx := context.WithValue(r.Context(), "session", session)
		r = r.WithContext(ctx)

		log.Printf("Authenticated user %v with session: %v", session.UserID, session.ID)
		next.ServeHTTP(w, r)
	})
}

func extractBearerToken(w http.ResponseWriter, r *http.Request) string {
	authHeaderValue := r.Header.Get("Authorization")
	if authHeaderValue == "" {
		http.Error(w, "Missing Bearer Token", http.StatusForbidden)
		return ""
	}

	result := strings.Split(authHeaderValue, "Bearer ")
	if len(result) != 2 {
		http.Error(w, "Invalid Bearer Token", http.StatusForbidden)
		return ""
	}

	return result[1]
}
