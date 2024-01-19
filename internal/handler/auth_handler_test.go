package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/absoluteyl/tasks-go/pkg/testutils"
)

func TestCreateAuthHandler(t *testing.T) {
	req := prepareCreateAuthRequest(t)

	rr := httptest.NewRecorder()
	CreateAuthHandler(rr, req)

	testutils.HttpStatusShouldBe(t, rr, http.StatusOK)

	response := parseMapResponse(t, rr)

	token, ok := response["token"]
	assert.True(t, ok, "Token field not found in response")

	_, ok = token.(string)
	assert.True(t, ok, "Token field is not a string")
}
func parseMapResponse(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}
	return response
}

func prepareCreateAuthRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("POST", "/auth", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
