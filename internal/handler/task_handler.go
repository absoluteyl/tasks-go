package handler

import (
	"encoding/json"
	"github.com/absoluteyl/tasks-go/internal/model"
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
	setContentType(w)

	var taskData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		setErrResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if _, ok := taskData["name"]; !ok {
		setErrResponse(w, http.StatusBadRequest, "Missing name attribute in request body")
		return
	}

	createdTaskID, err := h.taskService.CreateTask(taskData["name"].(string))
	if err != nil {
		setErrResponse(w, http.StatusInternalServerError, "Error creating task")
		return
	}

	newTask, err := h.taskService.GetTaskByID(createdTaskID)
	if err != nil {
		setErrResponse(w, http.StatusInternalServerError, "Error getting created task")
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
	tasks, err := h.taskService.GetTasks()
	if err != nil {
		http.Error(w, "Error getting tasks", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"result": tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var targetTask model.Task
	err := json.NewDecoder(r.Body).Decode(&targetTask)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	existingTask, err := h.taskService.GetTaskByID(targetTask.ID)
	if err != nil || existingTask.ID == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	err = h.taskService.UpdateTask(&targetTask)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"result": targetTask,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/task/")
	if id == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(id)
	existingTask, err := h.taskService.GetTaskByID(taskID)
	if err != nil || existingTask.ID == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	err = h.taskService.DeleteTask(taskID)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}
}

func setContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func jsonEncode(w http.ResponseWriter, response map[string]interface{}) {
	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func setErrResponse(w http.ResponseWriter, httpCode int, msg string) {
	w.WriteHeader(httpCode)
	response := map[string]interface{}{
		"result": msg,
	}
	jsonEncode(w, response)
}
