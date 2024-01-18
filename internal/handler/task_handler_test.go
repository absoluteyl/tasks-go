package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/absoluteyl/tasks-go/internal/model"
	"github.com/absoluteyl/tasks-go/internal/repository"
	"github.com/absoluteyl/tasks-go/internal/service"
	"github.com/absoluteyl/tasks-go/testutils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
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
	t := &testing.T{}
	setup(t)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(t *testing.T) {
	var err error
	testDB, err = testutils.ConnectDB()
	if err != nil {
		t.Fatal(err)
	}

	err = testutils.PrepareTaskTable(testDB)
	if err != nil {
		t.Fatal(err)
	}

	taskRepo = repository.NewTaskRepository(testDB)
	taskService = service.NewTaskServiceWithRepository(taskRepo)
	taskHandler = NewTaskHandler(taskService)
}

func teardown() {
	err := testutils.RemoveDB()
	if err != nil {
		fmt.Print(err)
	}
}

func TestTaskHandler(t *testing.T) {
	t.Run("CreateTaskHandler", testCreateTaskHandler)
	t.Run("GetTasksHandler", testGetTaskHandler)
	t.Run("UpdateTaskHandler", testUpdateTaskHandler)
	t.Run("DeleteTaskHandler", testDeleteTaskHandler)
}

func testCreateTaskHandler(t *testing.T) {
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
	assert.Equal(t, expectedHTTPStatus, rr.Code, "Handler returned wrong status code")

	var response map[string]map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	result, ok := response["result"]
	assert.Truef(t, ok, "Result field not found in response")

	expectedID := 1
	expectedName := taskData["name"].(string)
	expectedStatus := 0

	assert.Equal(t, expectedID, int(result["id"].(float64)), "Unexpected task id in response")
	assert.Equal(t, expectedName, result["name"], "Unexpected task name in response")
	assert.Equal(t, expectedStatus, int(result["status"].(float64)), "Unexpected task status in response")
}

func testGetTaskHandler(t *testing.T) {
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
	assert.Equal(t, expectedHTTPStatus, rr.Code, "Handler returned wrong status code")

	var response map[string][]map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	result, ok := response["result"]
	assert.True(t, ok, "Result field not found in response")

	expectedLength := len(tasksData)
	assert.Equal(t, expectedLength, len(result), "Unexpected task length in response")

	for i, task := range tasksData {

		expectedID := task.ID
		expectedName := task.Name
		expectedStatus := task.Status

		assert.Equal(t, expectedID, int(result[i]["id"].(float64)), "Unexpected task id in response")
		assert.Equal(t, expectedName, result[i]["name"], "Unexpected task name in response")
		assert.Equal(t, expectedStatus, int(result[i]["status"].(float64)), "Unexpected task status in response")

	}
}

func testUpdateTaskHandler(t *testing.T) {
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
	assert.Equal(t, expectedHTTPStatus, rr.Code, "Handler returned wrong status code")

	var response map[string]map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	result, ok := response["result"]
	assert.Truef(t, ok, "Result field not found in response")

	expectedID := taskData.ID
	expectedName := taskData.Name
	expectedStatus := taskData.Status
	assert.Equal(t, expectedID, int(result["id"].(float64)), "Unexpected task id in response")
	assert.Equal(t, expectedName, result["name"], "Unexpected task name in response")
	assert.Equal(t, expectedStatus, int(result["status"].(float64)), "Unexpected task status in response")
}

func testDeleteTaskHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/task/1", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	taskHandler.DeleteTaskHandler(rr, req)

	expectedHTTPStatus := http.StatusOK
	assert.Equal(t, expectedHTTPStatus, rr.Code, "Handler returned wrong status code")
}
