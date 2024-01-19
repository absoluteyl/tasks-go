package auth

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func GenerateToken(iat int64) (string, error) {
	token := SetSignMethod()

	PrepareClaims(token, iat)

	return Sign(token)
}

func SetSignMethod() *jwt.Token {
	token := jwt.New(jwt.SigningMethodHS256)
	return token
}

func PrepareClaims(token *jwt.Token, iat int64) {
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = iat
}

func Sign(token *jwt.Token) (string, error) {
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
}
