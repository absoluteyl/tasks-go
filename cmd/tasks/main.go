package main

import (
	"github.com/absoluteyl/tasks-go/internal/handler"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
)

func main() {
	mux := bone.New()
	mux.Post("/tasks/create", http.HandlerFunc(handler.CreateTaskHandler))

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":8080")
}
