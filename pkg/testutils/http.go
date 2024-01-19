package testutils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetupAuthorizationHeader(req *http.Request, authStr string) {
	req.Header.Add("Authorization", authStr)
}

func ParseResponse(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}
	return response
}

func HttpStatusShouldBe(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(t, expectedStatus, rr.Code, "Handler returned wrong status code")
}

func HttpResponseShouldBe(t *testing.T, rr *httptest.ResponseRecorder, expectedBody string) {
	assert.Equal(t, expectedBody, rr.Body.String(), "Handler returned unexpected body")
}

func HTTPBodyShouldHaveResultField(t *testing.T, response map[string]interface{}) (interface{}, bool) {
	result, ok := response["result"]
	assert.True(t, ok, "Result field not found in response")
	return result, ok
}
