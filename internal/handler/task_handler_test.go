package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/absoluteyl/tasks-go/internal/model"
	"github.com/absoluteyl/tasks-go/internal/repository"
	"github.com/absoluteyl/tasks-go/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testDB *sql.DB
var taskRepo *repository.TaskRepository
var taskService *service.TaskService
var taskHandler *TaskHandler

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Printf("Error opening test database: %v", err)
		os.Exit(1)
	}
	testDB = db

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		status INTEGER DEFAULT 0
	);`
	_, err = testDB.Exec(createTableSQL)
	if err != nil {
		fmt.Printf("Error creating tasks table: %v", err)
		os.Exit(1)
	}

	taskRepo = repository.NewTaskRepository(testDB)
	taskService = service.NewTaskServiceWithRepository(taskRepo)
	taskHandler = NewTaskHandler(taskService)
}

func teardown() {
	err := os.Remove("test.db")
	if err != nil {
		fmt.Printf("Error removing test database: %v", err)
	}
}

func TestCreateTaskHandler(t *testing.T) {
	taskData := map[string]interface{}{
		"name": "Eat Dinner",
	}

	taskJson, err := json.Marshal(taskData)
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/task", bytes.NewBuffer(taskJson))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	taskHandler.CreateTaskHandler(rr, req)

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
		},
	}

	req, err := http.NewRequest("Get", "/tasks", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	taskHandler.GetTasksHandler(rr, req)

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

func TestUpdateTaskHandler(t *testing.T) {
	taskData := model.Task{
		ID:     1,
		Name:   "Eat Lunch",
		Status: 1,
	}

	taskJson, err := json.Marshal(taskData)
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", "/task/", bytes.NewBuffer(taskJson))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	taskHandler.UpdateTaskHandler(rr, req)

	expectedHTTPStatus := http.StatusOK
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

	expectedName := taskData.Name
	if name, ok := result["name"].(string); !ok || name != expectedName {
		t.Errorf("Unexpected task name in response: got %v, want %v", name, expectedName)
	}

	expectedStatus := taskData.Status
	if status, ok := result["status"].(float64); !ok || int(status) != expectedStatus {
		t.Errorf("Unexpected task status in response: got %v, want %v", status, expectedStatus)
	}

	expectedID := taskData.ID
	if id, ok := result["id"].(float64); !ok || int(id) != expectedID {
		t.Errorf("Unexpected task id in response: got %v, want %v", id, expectedID)
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/task/1", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	taskHandler.DeleteTaskHandler(rr, req)

	expectedHTTPStatus := http.StatusOK
	if status := rr.Code; status != expectedHTTPStatus {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, expectedHTTPStatus)
	}
}
