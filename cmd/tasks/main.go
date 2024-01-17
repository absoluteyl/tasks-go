package main

import (
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
)

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func main() {
	mux := bone.New()
	mux.Post("/tasks/create", http.HandlerFunc(createTaskHandler))

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":8080")
}
