package service

import (
	"database/sql"
	"fmt"
	. "github.com/absoluteyl/tasks-go/pkg/testutils"
	"os"
	"testing"

	"github.com/absoluteyl/tasks-go/internal/model"
	"github.com/absoluteyl/tasks-go/internal/repository"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB
var taskRepo *repository.TaskRepository
var taskService *TaskService

var taskData = model.Task{
	ID:     1,
	Name:   "Eat Dinner",
	Status: 0,
}

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
	taskService = NewTaskServiceWithRepository(taskRepo)
}

func teardown() {
	err := RemoveDB()
	if err != nil {
		fmt.Print(err)
	}
}

func TestTaskService(t *testing.T) {
	t.Run("Create", testCreate)
	t.Run("GetList", testGetList)
	t.Run("Update", testUpdate)
	t.Run("Delete", testDelete)
}

func testCreate(t *testing.T) {
	taskName := taskData.Name

	taskID, err := taskService.CreateTask(taskName)
	assert.NoError(t, err)
	assert.NotZero(t, taskID)
}

func testGetList(t *testing.T) {
	tasks, err := taskService.GetTasks()
	assert.NoError(t, err)
	assert.NotZero(t, len(tasks))
	assert.Len(t, tasks, 1)

	assert.Equal(t, taskData.ID, tasks[0].ID)
	assert.Equal(t, taskData.Name, tasks[0].Name)
	assert.Equal(t, taskData.Status, tasks[0].Status)
}

func testUpdate(t *testing.T) {
	taskData.Name = "Eat Lunch"
	taskData.Status = 1

	err := taskService.UpdateTask(&taskData)
	assert.NoError(t, err)

	updatedTask, err := taskService.GetTaskByID(taskData.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedTask)

	assert.Equal(t, taskData.ID, updatedTask.ID)
	assert.Equal(t, taskData.Name, updatedTask.Name)
	assert.Equal(t, taskData.Status, updatedTask.Status)
}

func testDelete(t *testing.T) {
	err := taskService.DeleteTask(taskData.ID)
	assert.NoError(t, err)

	_, err = taskService.GetTaskByID(taskData.ID)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}
