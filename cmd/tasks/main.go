package main

import (
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
)

type Task struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
	ID     int    `json:"id"`
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	newTask := Task{
		Name:   "Eat Dinner",
		Status: 0,
		ID:     1,
	}

	response := map[string]interface{}{
		"result": newTask,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	mux := bone.New()
	mux.Post("/tasks/create", http.HandlerFunc(createTaskHandler))

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":8080")
}
