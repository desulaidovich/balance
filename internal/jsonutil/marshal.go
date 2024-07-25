package jsonutil

import (
	"encoding/json"
	"net/http"
)

// if err != nil {
// 	w.WriteHeader(http.StatusBadRequest)

// 	jsonError = &JsonMessage{
// 		Code:    http.StatusBadRequest,
// 		Message: `Параметр "money" должен быть числом`,
// 	}

// 	request, _ = json.Marshal(map[string]JsonMessage{
// 		"error": *jsonError,
// 	})

// 	w.Write(request)
// 	return
// }

type JsonMessage struct {
	Code    uint   `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func MarshalResponse(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(response)

	if err != nil {
		// TODO: заменить на логгер
		panic(err)
	}

	w.WriteHeader(status)
	w.Write(data)
}
