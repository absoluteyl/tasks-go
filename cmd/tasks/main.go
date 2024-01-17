package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"net/http"
)

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func main() {
	mux := bone.New()
	mux.Post("/tasks/create", http.HandlerFunc(createTaskHandler))

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":8080")
}
