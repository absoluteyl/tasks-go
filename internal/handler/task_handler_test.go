package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTaskHandler(t *testing.T) {
	taskData := map[string]interface{}{
		"name": "Eat Dinner",
	}

	taskJson, err := json.Marshal(taskData)
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/tasks/create", bytes.NewBuffer(taskJson))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	CreateTaskHandler(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	var response map[string]map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	result, ok := response["result"]
	if !ok {
		t.Error("Result field not found in response")
	}

	expectedName := taskData["name"].(string)
	if name, ok := result["name"].(string); !ok || name != expectedName {
		t.Errorf("Unexpected task name in response: got %v, want %v", name, expectedName)
	}

	expectedStatus := 0
	if status, ok := result["status"].(float64); !ok || int(status) != expectedStatus {
		t.Errorf("Unexpected task status in response: got %v, want %v", status, expectedStatus)
	}

	expectedID := 1
	if id, ok := result["id"].(float64); !ok || int(id) != expectedID {
		t.Errorf("Unexpected task id in response: got %v, want %v", id, expectedID)
	}
}

func TestGetTaskHandler(t *testing.T) {
	req, err := http.NewRequest("Get", "/tasks/1", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	GetTaskHandler(rr, req)

	expectedStatus := http.StatusOK
	if status := rr.Code; status != expectedStatus {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, expectedStatus)
	}
}
