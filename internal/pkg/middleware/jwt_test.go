package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestJWTMiddleware_NoHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	jwtHandler := JWTMiddleware(nextHandler)

	jwtHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	expected := "Authorization header is missing or not in 'Bearer {token}' format\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestJWTMiddleware_HeaderInvalid(t *testing.T) {
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "InvalidFormat")

	rr := httptest.NewRecorder()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	jwtHandler := JWTMiddleware(nextHandler)

	jwtHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	expected := "Authorization header is missing or not in 'Bearer {token}' format\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestJWTMiddleware_TokenInvalid(t *testing.T) {
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer InvalidToken")

	rr := httptest.NewRecorder()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	jwtHandler := JWTMiddleware(nextHandler)

	jwtHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	expected := "Invalid or expired token\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestJWTMiddleware_TokenOlderThan1Minute(t *testing.T) {
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

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	expected := "Token is older than 1 minute\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
