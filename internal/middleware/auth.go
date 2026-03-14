package middleware

import (
	"net/http"
	"os"
)

func AuthMiddleware() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			expectedToken := os.Getenv("ADMIN_API_KEY")
			if expectedToken == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer "+expectedToken {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next(w, r)
		}
	}
}
