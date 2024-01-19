package middleware

import (
	"github.com/absoluteyl/tasks-go/pkg/auth"
	. "github.com/absoluteyl/tasks-go/pkg/testutils"
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
	req := prepareGetTasksRequest(t)

	rr, jwtHandler := prepareHandlerRecorderWithMiddleware()
	jwtHandler.ServeHTTP(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusUnauthorized)

	response := ParseResponse(t, rr)
	ResultShouldExist(t, response)
	ResultShouldBe(t, ErrInvalidAuthorization, response["result"])
}

func testHeaderInvalid(t *testing.T) {
	req := prepareGetTasksRequest(t)
	SetupAuthorizationHeader(req, "InvalidFormat")

	rr, jwtHandler := prepareHandlerRecorderWithMiddleware()
	jwtHandler.ServeHTTP(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusUnauthorized)

	response := ParseResponse(t, rr)
	ResultShouldExist(t, response)
	ResultShouldBe(t, ErrInvalidAuthorization, response["result"])
}

func testTokenInvalid(t *testing.T) {
	req := prepareGetTasksRequest(t)
	SetupAuthorizationHeader(req, "Bearer InvalidToken")

	rr, jwtHandler := prepareHandlerRecorderWithMiddleware()
	jwtHandler.ServeHTTP(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusUnauthorized)
	response := ParseResponse(t, rr)
	ResultShouldExist(t, response)
	ResultShouldBe(t, ErrInvalidToken, response["result"])
}

func testTokenOlderThan1Minute(t *testing.T) {
	req := prepareGetTasksRequest(t)
	tokenString, err := auth.GenerateToken(time.Now().Add(-2 * time.Minute).Unix())
	assert.NoError(t, err)

	SetupAuthorizationHeader(req, "Bearer "+tokenString)

	rr, jwtHandler := prepareHandlerRecorderWithMiddleware()
	jwtHandler.ServeHTTP(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusUnauthorized)

	response := ParseResponse(t, rr)
	ResultShouldExist(t, response)
	ResultShouldBe(t, ErrTokenExpired, response["result"])
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
