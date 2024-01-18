package middleware

import (
	"net/http"
	"strings"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Authorization header is missing or not in 'Bearer {token}' format", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
