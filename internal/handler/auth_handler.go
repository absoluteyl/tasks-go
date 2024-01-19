package handler

import (
	"github.com/absoluteyl/tasks-go/pkg/auth"
	"net/http"
)

func CreateAuthHandler(w http.ResponseWriter, r *http.Request) {
	setContentType(w)

	// TODO: Authenticate if request is allow to gen token

	token, err := auth.GenerateToken()
	if err != nil {
		setErrResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := map[string]interface{}{
		"result": token,
	}
	jsonEncode(w, response)
}
