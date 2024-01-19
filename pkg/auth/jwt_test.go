package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	t.Run("GenerateToken", testGenerateToken)
	t.Run("ParseToken", testParseToken)
}

func testGenerateToken(t *testing.T) {
	tokenString, err := GenerateToken(time.Now().Unix())
	assert.NoError(t, err, "Failed to generate token")
	assert.NotEmpty(t, tokenString, "Generated token is empty")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	assert.NoError(t, err, "Failed to parse token")
	assert.True(t, token.Valid, "Token is not valid")
	assert.Equal(t, jwt.SigningMethodHS256, token.Method, "Invalid signing method")

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok, "Invalid token claims")

	iat, ok := claims["iat"].(float64)
	assert.True(t, ok, "Invalid issued at claim")
	assert.True(t, time.Now().After(time.Unix(int64(iat), 0)), "Token issued at time is in the future")
}

func testParseToken(t *testing.T) {
	tokenString, err := GenerateToken(time.Now().Unix())
	assert.NoError(t, err, "Failed to generate token")

	_, err = ParseToken(tokenString)
	assert.NoError(t, err, "Failed to parse token")
}
