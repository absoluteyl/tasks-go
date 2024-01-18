package handler

import (
	"encoding/json"
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
	if status := rr.Code; status != expectedHTTPStatus {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, expectedHTTPStatus)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	token, ok := response["token"]
	if !ok {
		t.Error("Result field not found in response")
	}

	_, ok = token.(string)
	if !ok {
		t.Error("Token field is not a string")
	}
}
