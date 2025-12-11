package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func JSONMessage(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Message: message})
}

func JSONError(w http.ResponseWriter, status int, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Error: err})
}
