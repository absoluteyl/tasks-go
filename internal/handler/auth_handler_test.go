package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAuthHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/auth", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	CreateAuthHandler(rr, req)

	expectedHTTPStatus := http.StatusOK
	assert.Equal(t, expectedHTTPStatus, rr.Code, "Handler returned wrong status code")

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	token, ok := response["token"]
	assert.Truef(t, ok, "Token field not found in response")

	_, ok = token.(string)
	assert.Truef(t, ok, "Token field is not a string")
}
