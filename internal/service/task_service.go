package service

import (
	"github.com/absoluteyl/tasks-go/internal/model"
	"github.com/absoluteyl/tasks-go/internal/repository"
)

type TaskService struct {
	taskRepository *repository.TaskRepository
}

func NewTaskServiceWithRepository(taskRepository *repository.TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}

func (s *TaskService) CreateTask(task *model.Task) (int, error) {
	return s.taskRepository.CreateTask(task)
}

func (s *TaskService) GetTasks() ([]model.Task, error) {
	tasks, err := s.taskRepository.GetTasks()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(task *model.Task) error {
	return s.taskRepository.UpdateTask(task)
}
