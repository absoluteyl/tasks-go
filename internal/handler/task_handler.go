package handler

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// FIXME: Temporary fixed status and id
	newTask.ID = 1
	newTask.Status = 0

	response := map[string]interface{}{
		"result": newTask,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME: Temporary fixed tasks
	tasks := []Task{
		{
			ID:     1,
			Name:   "Eat Dinner",
			Status: 0,
		},
		{
			ID:     2,
			Name:   "Go to sleep",
			Status: 0,
		},
	}

	response := map[string]interface{}{
		"result": tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
