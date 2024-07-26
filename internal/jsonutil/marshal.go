package jsonutil

import (
	"encoding/json"
	"net/http"
)

type JsonMessage struct {
	Code    uint   `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func MarshalResponse(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		message, _ := json.Marshal(&JsonMessage{
			Code:    http.StatusAlreadyReported,
			Message: err.Error(),
		})
		w.Write(message)
		return
	}

	w.WriteHeader(status)
	w.Write(data)
}
