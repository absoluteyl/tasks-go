package repository

import (
	"database/sql"
	"fmt"
	"github.com/absoluteyl/tasks-go/internal/model"
	. "github.com/absoluteyl/tasks-go/pkg/testutils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testDB *sql.DB
var taskRepo *TaskRepository
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

	taskRepo = NewTaskRepository(testDB)
}

func teardown() {
	err := RemoveDB()
	if err != nil {
		fmt.Print(err)
	}
}

func TestTaskRepository(t *testing.T) {
	t.Run("Create", testCreate)
	t.Run("GetList", testGetList)
	t.Run("Update", testUpdate)
	t.Run("Delete", testDelete)
}

func testCreate(t *testing.T) {
	taskName := taskData.Name

	taskID, err := taskRepo.CreateTask(taskName)
	assert.NoError(t, err)
	assert.NotZero(t, taskID)
}

func testGetList(t *testing.T) {
	tasks, err := taskRepo.GetTasks()
	assert.NoError(t, err)
	assert.NotEmpty(t, tasks)
	assert.Len(t, tasks, 1)

	assert.Equal(t, taskData.ID, tasks[0].ID)
	assert.Equal(t, taskData.Name, tasks[0].Name)
	assert.Equal(t, taskData.Status, tasks[0].Status)
}

func testUpdate(t *testing.T) {
	taskData.Name = "Eat Lunch"
	taskData.Status = 1

	err := taskRepo.UpdateTask(&taskData)
	assert.NoError(t, err)

	updatedTask, err := taskRepo.GetTaskByID(taskData.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedTask)

	assert.Equal(t, taskData.ID, updatedTask.ID)
	assert.Equal(t, taskData.Name, updatedTask.Name)
	assert.Equal(t, taskData.Status, updatedTask.Status)
}

func testDelete(t *testing.T) {
	err := taskRepo.DeleteTask(taskData.ID)
	assert.NoError(t, err)

	_, err = taskRepo.GetTaskByID(taskData.ID)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}
