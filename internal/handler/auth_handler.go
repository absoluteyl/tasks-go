package handler

import (
	"github.com/absoluteyl/tasks-go/pkg/auth"
	"net/http"
	"os"
)

var (
	// FIXME: Get user credentials somewhere else
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
)

func CreateAuthHandler(w http.ResponseWriter, r *http.Request) {
	SetContentType(w)

	user, pass, ok := r.BasicAuth()
	if !ok || user != username || pass != password {
		w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
		SetErrResponse(w, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		SetErrResponse(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	response := map[string]interface{}{
		"result": token,
	}
	jsonEncode(w, response)
}
