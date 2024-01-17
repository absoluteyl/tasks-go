package main

import (
	"database/sql"
	"github.com/absoluteyl/tasks-go/internal/handler"
	"github.com/absoluteyl/tasks-go/internal/repository"
	"github.com/absoluteyl/tasks-go/internal/service"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/tasks.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	taskRepository := repository.NewTaskRepository(db)
	taskService := service.NewTaskServiceWithRepository(taskRepository)
	taskHandler := handler.NewTaskHandler(taskService)

	mux := bone.New()
	mux.Post("/tasks", http.HandlerFunc(taskHandler.CreateTaskHandler))
	mux.Get("/tasks", http.HandlerFunc(taskHandler.GetTasksHandler))

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":8080")
}
