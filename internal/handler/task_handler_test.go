package handler

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/absoluteyl/tasks-go/internal/model"
	"github.com/absoluteyl/tasks-go/internal/repository"
	"github.com/absoluteyl/tasks-go/internal/service"
	. "github.com/absoluteyl/tasks-go/pkg/testutils"
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
	testDB, err = ConnectDB()
	if err != nil {
		t.Fatal(err)
	}

	err = PrepareTaskTable(testDB)
	if err != nil {
		t.Fatal(err)
	}

	taskRepo = repository.NewTaskRepository(testDB)
	taskService = service.NewTaskServiceWithRepository(taskRepo)
	taskHandler = NewTaskHandler(taskService)
}

func teardown() {
	err := RemoveDB()
	if err != nil {
		fmt.Print(err)
	}
}

func TestTaskHandler(t *testing.T) {
	t.Run("CreateMissingName", testCreateMissingName)
	t.Run("Create", testCreate)

	t.Run("GetList", testGetList)

	t.Run("UpdateWithoutID", testUpdateWithoutID)
	t.Run("UpdateWithoutID", testUpdateWithInvalidID)
	t.Run("UpdateWithIDInBody", testUpdateWithIDInBody)
	t.Run("UpdateNotExist", testUpdateNotExist)
	t.Run("Update", testUpdate)

	t.Run("DeleteNotExist", testDeleteNotExist)
	t.Run("Delete", testDelete)
}

func testCreateMissingName(t *testing.T) {
	taskData := map[string]interface{}{}

	reqBody := PrepareJsonBody(t, taskData)

	req := prepareCreateTaskRequest(t, reqBody)

	rr := httptest.NewRecorder()
	taskHandler.CreateTaskHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusBadRequest)

	response := ParseResponse(t, rr)
	HTTPBodyShouldHaveResultField(t, response)
}

func testCreate(t *testing.T) {
	taskData := map[string]interface{}{
		"name": "Eat Dinner",
	}

	reqBody := PrepareJsonBody(t, taskData)

	req := prepareCreateTaskRequest(t, reqBody)

	rr := httptest.NewRecorder()
	taskHandler.CreateTaskHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusCreated)

	response := ParseMapMapResponse(t, rr)

	result, ok := response["result"]
	assert.True(t, ok, "Result field not found in response")

	taskShouldBe(t, model.Task{
		ID:     1,
		Name:   taskData["name"].(string),
		Status: 0,
	}, model.Task{
		ID:     int(result["id"].(float64)),
		Name:   result["name"].(string),
		Status: int(result["status"].(float64)),
	})
}

func testGetList(t *testing.T) {
	tasksData := []model.Task{
		{
			ID:     1,
			Name:   "Eat Dinner",
			Status: 0,
		},
	}

	req := prepareGetTasksRequest(t)

	rr := httptest.NewRecorder()
	taskHandler.GetTasksHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusOK)

	response := ParseMapArrayResponse(t, rr)

	result, ok := response["result"]
	assert.True(t, ok, "Result field not found in response")

	expectedLength := len(tasksData)
	assert.Equal(t, expectedLength, len(result), "Unexpected task length in response")

	for i, task := range tasksData {
		taskShouldBe(t, task, model.Task{
			ID:     int(result[i]["id"].(float64)),
			Name:   result[i]["name"].(string),
			Status: int(result[i]["status"].(float64)),
		})
	}
}

func testUpdateWithoutID(t *testing.T) {
	taskData := map[string]interface{}{
		"name":   "Eat Lunch",
		"status": 1,
	}

	reqBody := PrepareJsonBody(t, taskData)

	req, err := http.NewRequest("PUT", "/task/", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	taskHandler.UpdateTaskHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusBadRequest)

	response := ParseResponse(t, rr)
	HTTPBodyShouldHaveResultField(t, response)
	ResultShouldBe(t, "Missing task ID", response["result"])
}

func testUpdateWithInvalidID(t *testing.T) {
	taskData := map[string]interface{}{
		"name":   "Eat Lunch",
		"status": 1,
	}

	reqBody := PrepareJsonBody(t, taskData)

	req, err := http.NewRequest("PUT", "/task/invalid", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	taskHandler.UpdateTaskHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusBadRequest)

	response := ParseResponse(t, rr)
	HTTPBodyShouldHaveResultField(t, response)
	ResultShouldBe(t, "Invalid task ID", response["result"])
}

func testUpdateWithIDInBody(t *testing.T) {
	taskData := map[string]interface{}{
		"id": 1,
	}

	reqBody := PrepareJsonBody(t, taskData)

	req := prepareUpdateTaskRequest(t, taskData["id"].(int), reqBody)

	rr := httptest.NewRecorder()
	taskHandler.UpdateTaskHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusBadRequest)

	response := ParseResponse(t, rr)
	HTTPBodyShouldHaveResultField(t, response)
	ResultShouldBe(t, "Invalid request body: task ID not allowed", response["result"])
}

func testUpdateNotExist(t *testing.T) {
	taskData := map[string]interface{}{
		"name":   "Eat Lunch",
		"status": 1,
	}

	reqBody := PrepareJsonBody(t, taskData)

	req := prepareUpdateTaskRequest(t, 2, reqBody)

	rr := httptest.NewRecorder()
	taskHandler.UpdateTaskHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusNotFound)
}

func testUpdate(t *testing.T) {
	taskData := map[string]interface{}{
		"name":   "Eat Lunch",
		"status": 1,
	}

	reqBody := PrepareJsonBody(t, taskData)

	req := prepareUpdateTaskRequest(t, 1, reqBody)

	rr := httptest.NewRecorder()
	taskHandler.UpdateTaskHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusOK)

	response := ParseMapMapResponse(t, rr)

	result, ok := response["result"]
	assert.True(t, ok, "Result field not found in response")

	taskShouldBe(t, model.Task{
		ID:     1,
		Name:   taskData["name"].(string),
		Status: taskData["status"].(int),
	}, model.Task{
		ID:     int(result["id"].(float64)),
		Name:   result["name"].(string),
		Status: int(result["status"].(float64)),
	})
}

func testDeleteNotExist(t *testing.T) {
	req := prepareDeleteTaskRequest(t, 999)

	rr := httptest.NewRecorder()
	taskHandler.DeleteTaskHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusNotFound)
}

func testDelete(t *testing.T) {
	req := prepareDeleteTaskRequest(t, 1)

	rr := httptest.NewRecorder()
	taskHandler.DeleteTaskHandler(rr, req)

	HttpStatusShouldBe(t, rr, http.StatusOK)
}

func prepareCreateTaskRequest(t *testing.T, body []byte) *http.Request {
	req, err := http.NewRequest("POST", "/task", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	return req
}

func prepareGetTasksRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	return req
}

func prepareUpdateTaskRequest(t *testing.T, id int, body []byte) *http.Request {
	taskID := fmt.Sprintf("%d", id)
	req, err := http.NewRequest("PUT", "/task/"+taskID, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	return req
}

func prepareDeleteTaskRequest(t *testing.T, id int) *http.Request {
	taskID := fmt.Sprintf("%d", id)
	req, err := http.NewRequest("DELETE", "/task/"+taskID, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	return req
}

func taskShouldBe(t *testing.T, expectedTask model.Task, actualTask model.Task) {
	assert.Equal(t, expectedTask.ID, actualTask.ID, "Unexpected task id")
	assert.Equal(t, expectedTask.Name, actualTask.Name, "Unexpected task name")
	assert.Equal(t, expectedTask.Status, actualTask.Status, "Unexpected task status")
}
