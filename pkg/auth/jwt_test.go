package auth

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	tokenString, err := GenerateToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if tokenString == "" {
		t.Fatal("Generated token is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("Unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}

		return []byte("secret"), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Fatal("Token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Invalid token claims")
	}

	if iat, ok := claims["iat"].(float64); ok {
		if time.Now().Before(time.Unix(int64(iat), 0)) {
			t.Fatal("Token issued at time is in the future")
		}
	} else {
		t.Fatal("Invalid issued at claim")
	}
}

func TestParseToken(t *testing.T) {
	tokenString, err := GenerateToken()
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	_, err = ParseToken(tokenString)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}
}
