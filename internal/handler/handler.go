package handler

import (
	"encoding/json"
	"net/http"
)

func SetContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SetErrResponse(w http.ResponseWriter, httpCode int, msg string) {
	w.WriteHeader(httpCode)
	response := map[string]interface{}{
		"result": msg,
	}
	jsonEncode(w, response)
}

func jsonEncode(w http.ResponseWriter, response map[string]interface{}) {
	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
