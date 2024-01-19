package repository

import (
	"database/sql"
	"fmt"
	"github.com/absoluteyl/tasks-go/internal/model"
	"github.com/absoluteyl/tasks-go/pkg/testutils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testDB *sql.DB
var taskRepo *TaskRepository

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

	taskRepo = NewTaskRepository(testDB)
}

func teardown() {
	err := testutils.RemoveDB()
	if err != nil {
		fmt.Print(err)
	}
}

func TestTaskRepository(t *testing.T) {
	t.Run("CreateTask", testCreateTask)
	t.Run("GetTasks", testGetTasks)
	t.Run("UpdateTask", testUpdateTask)
	t.Run("DeleteTask", testDeleteTask)
}

func testCreateTask(t *testing.T) {
	taskName := "Eat Dinner"

	taskID, err := taskRepo.CreateTask(taskName)
	assert.NoError(t, err)
	assert.NotZero(t, taskID)
}

func testGetTasks(t *testing.T) {
	tasks, err := taskRepo.GetTasks()
	assert.NoError(t, err)
	assert.NotEmpty(t, tasks)
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

	err := taskRepo.UpdateTask(task)
	assert.NoError(t, err)

	tasks, err := taskRepo.GetTasks()
	assert.NoError(t, err)
	assert.NotEmpty(t, tasks)
	assert.Len(t, tasks, 1)

	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "Eat Lunch", tasks[0].Name)
	assert.Equal(t, 1, tasks[0].Status)
}

func testDeleteTask(t *testing.T) {
	err := taskRepo.DeleteTask(1)
	assert.NoError(t, err)

	tasks, err := taskRepo.GetTasks()
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}
