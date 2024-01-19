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

func (s *TaskService) CreateTask(taskName string) (int, error) {
	return s.taskRepository.CreateTask(taskName)
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

func (s *TaskService) DeleteTask(id int) error {
	return s.taskRepository.DeleteTask(id)
}

func (s *TaskService) GetTaskByID(id int) (model.Task, error) {
	return s.taskRepository.GetTaskByID(id)
}
