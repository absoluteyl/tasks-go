package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/absoluteyl/tasks-go/pkg/testutils"
)

func TestAuthHandler(t *testing.T) {
	t.Run("Create", testCreateAuthWithInvalidCredentials)
	t.Run("Create", testCreateAuth)
}

func testCreateAuthWithInvalidCredentials(t *testing.T) {
	req := prepareCreateAuthRequest(t, "invalid", "invalid")

	rr := httptest.NewRecorder()
	CreateAuthHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusUnauthorized)

	response := ParseResponse(t, rr)
	ResultShouldBe(t, "Unauthorized", response["result"])
}

func testCreateAuth(t *testing.T) {
	req := prepareCreateAuthRequest(t, username, password)

	rr := httptest.NewRecorder()
	CreateAuthHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusOK)

	response := ParseResponse(t, rr)
	result, ok := HTTPBodyShouldHaveResultField(t, response)

	_, ok = result.(string)
	assert.True(t, ok, "Result field is not a string")
}

func prepareCreateAuthRequest(t *testing.T, user string, pass string) *http.Request {
	req, err := http.NewRequest("POST", "/auth", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth(user, pass)
	return req
}
