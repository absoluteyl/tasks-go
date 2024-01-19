package main

import (
	"database/sql"
	"github.com/absoluteyl/tasks-go/internal/handler"
	"github.com/absoluteyl/tasks-go/internal/middleware"
	"github.com/absoluteyl/tasks-go/internal/repository"
	"github.com/absoluteyl/tasks-go/internal/service"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := connectDB("sqlite3", "./db/tasks.db")
	defer db.Close()

	mux := setupRouter(db)

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":8080")
}

func setupRouter(db *sql.DB) *bone.Mux {
	taskRepository := repository.NewTaskRepository(db)
	taskService := service.NewTaskServiceWithRepository(taskRepository)
	taskHandler := handler.NewTaskHandler(taskService)

	mux := bone.New()

	mux.Post("/auth", http.HandlerFunc(handler.CreateAuthHandler))

	mux.Post("/task", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.CreateTaskHandler)))
	mux.Get("/tasks", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.GetTasksHandler)))
	mux.Put("/task/:id", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.UpdateTaskHandler)))
	mux.Delete("/task/:id", middleware.JWTMiddleware(http.HandlerFunc(taskHandler.DeleteTaskHandler)))

	return mux
}

func connectDB(dbDriver string, dbPath string) *sql.DB {
	db, err := sql.Open(dbDriver, dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	return db
}
