package middleware

import (
	. "github.com/absoluteyl/tasks-go/internal/handler"
	"github.com/absoluteyl/tasks-go/pkg/auth"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetContentType(w)

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			SetErrResponse(w, http.StatusUnauthorized, ErrInvalidAuthorization)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := auth.ParseToken(tokenString)
		if err != nil || !token.Valid {
			SetErrResponse(w, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			SetErrResponse(w, http.StatusUnauthorized, ErrInvalidClaims)
			return
		}

		if iat, ok := claims["iat"].(float64); ok {
			if time.Now().After(time.Unix(int64(iat), 0).Add(1 * time.Minute)) {
				SetErrResponse(w, http.StatusUnauthorized, ErrTokenExpired)
				return
			}
		} else {
			SetErrResponse(w, http.StatusUnauthorized, ErrInvalidIssueAt)
			return
		}

		next.ServeHTTP(w, r)
	})
}
