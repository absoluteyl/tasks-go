package repository

import (
	"database/sql"
	"fmt"
	"github.com/absoluteyl/tasks-go/internal/model"
	"github.com/absoluteyl/tasks-go/testutils"
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
	err := os.Remove("test.db")
	if err != nil {
		fmt.Printf("Error removing test database: %v", err)
	}
}

func TestCreateTask(t *testing.T) {
	task := &model.Task{
		Name: "Eat Dinner",
	}

	taskID, err := taskRepo.CreateTask(task)
	assert.NoError(t, err)
	assert.NotZero(t, taskID)
}

func TestGetTasks(t *testing.T) {
	tasks, err := taskRepo.GetTasks()
	assert.NoError(t, err)
	assert.NotEmpty(t, tasks)
	assert.Len(t, tasks, 1)

	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "Eat Dinner", tasks[0].Name)
	assert.Equal(t, 0, tasks[0].Status)
}

func TestUpdateTask(t *testing.T) {
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

func TestDeleteTask(t *testing.T) {
	err := taskRepo.DeleteTask(1)
	assert.NoError(t, err)

	tasks, err := taskRepo.GetTasks()
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}
