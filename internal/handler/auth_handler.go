package handler

import (
	"encoding/json"
	"github.com/absoluteyl/tasks-go/pkg/auth"
	"net/http"
)

func CreateAuthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Authenticate if request is allow to gen token

	token, err := auth.GenerateToken()
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"result": token,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
