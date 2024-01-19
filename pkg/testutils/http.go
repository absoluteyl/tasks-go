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

func PrepareJsonBody(t *testing.T, data map[string]interface{}) []byte {
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}
	return body
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

func ResultShouldBe(t *testing.T, expected interface{}, result interface{}) {
	assert.Equal(t, expected, result, "Result field value is unexpected")
}

func ParseMapArrayResponse(t *testing.T, rr *httptest.ResponseRecorder) map[string][]map[string]interface{} {
	var response map[string][]map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON response: %v", err)
	}
	return response
}

func ParseMapMapResponse(t *testing.T, rr *httptest.ResponseRecorder) map[string]map[string]interface{} {
	var response map[string]map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}
	return response
}
