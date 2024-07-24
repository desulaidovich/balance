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
	Message string `json:"message"`
	Node    any    `json:"response"`
}

func MarshalResponse(w http.ResponseWriter, status int, name string, response any) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(response)

	if err != nil {
		panic(err)
	}

	w.WriteHeader(status)
	w.Write(data)
}
