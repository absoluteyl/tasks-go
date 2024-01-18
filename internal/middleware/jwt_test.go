package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestJWTMiddleware(t *testing.T) {
	t.Run("NoHeader", testNoHeader)
	t.Run("HeaderInvalid", testHeaderInvalid)
	t.Run("TokenInvalid", testTokenInvalid)
	t.Run("TokenOlderThan1Minute", testTokenOlderThan1Minute)
}

func testNoHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	jwtHandler := JWTMiddleware(nextHandler)

	jwtHandler.ServeHTTP(rr, req)

	expectedHTTPStatus := http.StatusUnauthorized
	assert.Equal(t, expectedHTTPStatus, rr.Code, "Handler returned wrong status code")

	expected := "Authorization header is missing or not in 'Bearer {token}' format\n"
	assert.Equal(t, expected, rr.Body.String(), "Handler returned unexpected body")
}

func testHeaderInvalid(t *testing.T) {
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "InvalidFormat")

	rr := httptest.NewRecorder()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	jwtHandler := JWTMiddleware(nextHandler)

	jwtHandler.ServeHTTP(rr, req)

	expectedHTTPStatus := http.StatusUnauthorized
	assert.Equal(t, expectedHTTPStatus, rr.Code, "Handler returned wrong status code")

	expected := "Authorization header is missing or not in 'Bearer {token}' format\n"
	assert.Equal(t, expected, rr.Body.String(), "Handler returned unexpected body")
}

func testTokenInvalid(t *testing.T) {
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer InvalidToken")

	rr := httptest.NewRecorder()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	jwtHandler := JWTMiddleware(nextHandler)

	jwtHandler.ServeHTTP(rr, req)

	expectedHTTPStatus := http.StatusUnauthorized
	assert.Equal(t, expectedHTTPStatus, rr.Code, "Handler returned wrong status code")

	expected := "Invalid or expired token\n"
	assert.Equal(t, expected, rr.Body.String(), "Handler returned unexpected body")
}

func testTokenOlderThan1Minute(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().Add(-2 * time.Minute).Unix()
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	jwtHandler := JWTMiddleware(nextHandler)

	jwtHandler.ServeHTTP(rr, req)

	expectedHTTPStatus := http.StatusUnauthorized
	assert.Equal(t, expectedHTTPStatus, rr.Code, "Handler returned wrong status code")

	expected := "Token is older than 1 minute\n"
	assert.Equal(t, expected, rr.Body.String(), "Handler returned unexpected body")

}
