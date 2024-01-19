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
	t.Run("CreateTask", testCreateTask)
	t.Run("GetTasks", testGetTasks)
	t.Run("UpdateTask", testUpdateTask)
	t.Run("DeleteTask", testDeleteTask)
}

func testCreateTask(t *testing.T) {
	taskName := "Eat Dinner"

	taskID, err := taskService.CreateTask(taskName)
	assert.NoError(t, err)
	assert.NotZero(t, taskID)
}

func testGetTasks(t *testing.T) {
	tasks, err := taskService.GetTasks()
	assert.NoError(t, err)
	assert.NotZero(t, len(tasks))
	assert.Len(t, tasks, 1)

	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "Eat Dinner", tasks[0].Name)
	assert.Equal(t, 0, tasks[0].Status)
}

func testUpdateTask(t *testing.T) {
	task := &model.Task{
		ID:     1,
		Name:   "Eat Lunch",
		Status: 1,
	}

	err := taskService.UpdateTask(task)
	assert.NoError(t, err)

	tasks, err := taskService.GetTasks()
	assert.NoError(t, err)
	assert.NotZero(t, len(tasks))
	assert.Len(t, tasks, 1)

	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "Eat Lunch", tasks[0].Name)
	assert.Equal(t, 1, tasks[0].Status)
}

func testDeleteTask(t *testing.T) {
	err := taskService.DeleteTask(1)
	assert.NoError(t, err)

	tasks, err := taskService.GetTasks()
	assert.NoError(t, err)
	assert.Zero(t, len(tasks))
}
