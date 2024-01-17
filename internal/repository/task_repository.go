package repository

import (
	"database/sql"
	"github.com/absoluteyl/tasks-go/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(task *model.Task) (int, error) {
	createTaskSQL := `
	INSERT INTO tasks (name) VALUES (?)
	`
	result, err := r.db.Exec(createTaskSQL, task.Name)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInsertID), nil
}
