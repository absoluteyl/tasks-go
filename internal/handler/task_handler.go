package handler

import (
	"encoding/json"
	"github.com/absoluteyl/tasks-go/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	SetContentType(w)

	var taskData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		SetErrResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	if _, ok := taskData["name"]; !ok {
		SetErrResponse(w, http.StatusBadRequest, ErrMissingTaskName)
		return
	}

	createdTaskID, err := h.taskService.CreateTask(taskData["name"].(string))
	if err != nil {
		SetErrResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	newTask, err := h.taskService.GetTaskByID(createdTaskID)
	if err != nil {
		SetErrResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"result": map[string]interface{}{
			"id":     newTask.ID,
			"name":   newTask.Name,
			"status": newTask.Status,
		},
	}
	jsonEncode(w, response)
}

func (h *TaskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	SetContentType(w)

	tasks, err := h.taskService.GetTasks()
	if err != nil {
		SetErrResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	response := map[string]interface{}{
		"result": tasks,
	}
	jsonEncode(w, response)
}

func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	SetContentType(w)

	id := strings.TrimPrefix(r.URL.Path, "/task/")
	if id == "" {
		SetErrResponse(w, http.StatusBadRequest, ErrMissingTaskID)
		return
	}

	taskID, err := strconv.Atoi(id)
	if err != nil {
		SetErrResponse(w, http.StatusBadRequest, ErrInvalidTaskID)
		return
	}

	var taskData map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		SetErrResponse(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	if taskData["id"] != nil {
		SetErrResponse(w, http.StatusBadRequest, ErrNotAllowTaskID)
		return
	}

	existingTask, err := h.taskService.GetTaskByID(taskID)
	if err != nil || existingTask.ID == 0 {
		SetErrResponse(w, http.StatusNotFound, ErrTaskNotFound)
		return
	}

	if taskData["name"] != nil {
		existingTask.Name = taskData["name"].(string)
	}

	if taskData["status"] != nil {
		existingTask.Status = int(taskData["status"].(float64))
	}

	err = h.taskService.UpdateTask(&existingTask)
	if err != nil {
		SetErrResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	response := map[string]interface{}{
		"result": existingTask,
	}
	jsonEncode(w, response)
}

func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	SetContentType(w)

	id := strings.TrimPrefix(r.URL.Path, "/task/")
	if id == "" {
		SetErrResponse(w, http.StatusBadRequest, ErrMissingTaskID)
		return
	}

	taskID, err := strconv.Atoi(id)
	if err != nil {
		SetErrResponse(w, http.StatusBadRequest, ErrInvalidTaskID)
		return
	}

	existingTask, err := h.taskService.GetTaskByID(taskID)
	if err != nil || existingTask.ID == 0 {
		SetErrResponse(w, http.StatusNotFound, ErrTaskNotFound)
		return
	}

	err = h.taskService.DeleteTask(taskID)
	if err != nil {
		SetErrResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}
}
