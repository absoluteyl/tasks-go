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
