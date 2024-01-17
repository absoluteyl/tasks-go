package handler

import (
	"bytes"
	"encoding/json"
	"github.com/absoluteyl/tasks-go/internal/model"
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

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskJson))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	CreateTaskHandler(rr, req)

	expectedHTTPStatus := http.StatusCreated
	if status := rr.Code; status != expectedHTTPStatus {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, expectedHTTPStatus)
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
	tasksData := []model.Task{
		{
			ID:     1,
			Name:   "Eat Dinner",
			Status: 0,
		}, {
			ID:     2,
			Name:   "Go to sleep",
			Status: 0,
		},
	}

	req, err := http.NewRequest("Get", "/tasks", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	GetTasksHandler(rr, req)

	expectedHTTPStatus := http.StatusOK
	if status := rr.Code; status != expectedHTTPStatus {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, expectedHTTPStatus)
	}

	var response map[string][]map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	result, ok := response["result"]
	if !ok {
		t.Error("Result field not found in response")
	}

	expectedLength := len(tasksData)
	if len(result) != expectedLength {
		t.Errorf("Unexpected task length in response: got %v, want %v", len(result), expectedLength)
	}

	for i, task := range tasksData {
		expectedName := task.Name
		if name, ok := result[i]["name"].(string); !ok || name != expectedName {
			t.Errorf("Unexpected task name in response: got %v, want %v", name, expectedName)
		}

		expectedStatus := task.Status
		if status, ok := result[i]["status"].(float64); !ok || int(status) != expectedStatus {
			t.Errorf("Unexpected task status in response: got %v, want %v", status, expectedStatus)
		}

		expectedID := task.ID
		if id, ok := result[i]["id"].(float64); !ok || int(id) != expectedID {
			t.Errorf("Unexpected task id in response: got %v, want %v", id, expectedID)
		}
	}
}
