package middlewares

import (
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

		err := authService.ValidateAuthToken(r.Context(), authToken)
		if err != nil {
			http.Error(w, "Invalid Bearer Token", http.StatusForbidden)
			return
		}

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
