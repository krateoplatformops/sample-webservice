package auth

import (
	"net/http"
	"strings"
)

var (
	serviceName = "chart-inspector"
	// validToken  = "your-secret-token" // Replace with your actual Bearer token
)

// Middleware for Bearer authentication
func BearerAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		spl := strings.Split(authHeader, "Bearer ")

		if len(spl) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
