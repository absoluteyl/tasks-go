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

func (r *TaskRepository) GetTasks() ([]model.Task, error) {
	getTasksSQL := `
	SELECT id, name, status FROM tasks
	`
	rows, err := r.db.Query(getTasksSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Name, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(task *model.Task) error {
	updateTaskSQL := `
	UPDATE tasks SET name = ?, status = ? WHERE id = ?
	`
	_, err := r.db.Exec(updateTaskSQL, task.Name, task.Status, task.ID)
	if err != nil {
		return err
	}

	return nil
}
