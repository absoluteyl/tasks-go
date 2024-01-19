package middleware

import (
	"github.com/absoluteyl/tasks-go/pkg/auth"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, ErrInvalidAuthorization, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := auth.ParseToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, ErrInvalidClaims, http.StatusUnauthorized)
			return
		}

		if iat, ok := claims["iat"].(float64); ok {
			if time.Now().After(time.Unix(int64(iat), 0).Add(1 * time.Minute)) {
				http.Error(w, ErrTokenExpired, http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, ErrInvalidIssueAt, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
