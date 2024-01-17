package service

import (
	"database/sql"
	"fmt"
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
	taskService = NewTaskServiceWithRepository(taskRepo)
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

	taskID, err := taskService.CreateTask(task)
	assert.NoError(t, err)
	assert.NotZero(t, taskID)
}

func TestGetTasks(t *testing.T) {
	tasks, err := taskService.GetTasks()
	assert.NoError(t, err)
	assert.NotZero(t, len(tasks))
	assert.Len(t, tasks, 1)

	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "Eat Dinner", tasks[0].Name)
	assert.Equal(t, 0, tasks[0].Status)
}
