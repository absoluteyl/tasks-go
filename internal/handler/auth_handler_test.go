package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/absoluteyl/tasks-go/pkg/testutils"
)

func TestCreateAuthHandler(t *testing.T) {
	req := prepareCreateAuthRequest(t)

	rr := httptest.NewRecorder()
	CreateAuthHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusOK)

	response := ParseResponse(t, rr)
	result, ok := HTTPBodyShouldHaveResultField(t, response)

	_, ok = result.(string)
	assert.True(t, ok, "Result field is not a string")
}

func prepareCreateAuthRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/auth", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
