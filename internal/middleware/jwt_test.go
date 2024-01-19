package middleware

import (
	"github.com/absoluteyl/tasks-go/pkg/testutils"
	"github.com/dgrijalva/jwt-go"
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
	req := prepareGetTasksRequest(t)

	rr, jwtHandler := prepareHandlerRecorderWithMiddleware()
	jwtHandler.ServeHTTP(rr, req)

	testutils.HttpStatusShouldBe(t, rr, http.StatusUnauthorized)
	testutils.HttpResponseShouldBe(t, rr, "Authorization header is missing or not in 'Bearer {token}' format\n")
}

func testHeaderInvalid(t *testing.T) {
	req := prepareGetTasksRequest(t)
	testutils.SetupAuthorizationHeader(req, "InvalidFormat")

	rr, jwtHandler := prepareHandlerRecorderWithMiddleware()
	jwtHandler.ServeHTTP(rr, req)

	testutils.HttpStatusShouldBe(t, rr, http.StatusUnauthorized)
	testutils.HttpResponseShouldBe(t, rr, "Authorization header is missing or not in 'Bearer {token}' format\n")
}

func testTokenInvalid(t *testing.T) {
	req := prepareGetTasksRequest(t)
	testutils.SetupAuthorizationHeader(req, "Bearer InvalidToken")

	rr, jwtHandler := prepareHandlerRecorderWithMiddleware()
	jwtHandler.ServeHTTP(rr, req)

	testutils.HttpStatusShouldBe(t, rr, http.StatusUnauthorized)
	testutils.HttpResponseShouldBe(t, rr, "Invalid or expired token\n")
}

func testTokenOlderThan1Minute(t *testing.T) {
	req := prepareGetTasksRequest(t)
	tokenString := generateJWTToken(t)

	testutils.SetupAuthorizationHeader(req, "Bearer "+tokenString)

	rr, jwtHandler := prepareHandlerRecorderWithMiddleware()
	jwtHandler.ServeHTTP(rr, req)

	testutils.HttpStatusShouldBe(t, rr, http.StatusUnauthorized)
	testutils.HttpResponseShouldBe(t, rr, "Token is older than 1 minute\n")
}

func prepareGetTasksRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	return req
}

func prepareHandlerRecorderWithMiddleware() (*httptest.ResponseRecorder, http.Handler) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middlewareHandler := JWTMiddleware(handler)
	return rr, middlewareHandler
}

func generateJWTToken(t *testing.T) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().Add(-2 * time.Minute).Unix()
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	return tokenString
}
