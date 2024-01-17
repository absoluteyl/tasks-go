package service

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/absoluteyl/tasks-go/internal/repository"
	_ "github.com/mattn/go-sqlite3"
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
	t.Skip("Skipping test for now")
}
