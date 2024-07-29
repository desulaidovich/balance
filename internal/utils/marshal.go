package utils

import (
	"encoding/json"
	"net/http"
)

const (
	REQUEST_ERROR_CODE    = 500
	REQUEST_NO_ERROR_CODE = 0
)

type JSONMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    *Data  `json:"data,omitempty"`
}

type Data struct {
	WalletData *WalletData `json:"wallet_data,omitempty"`
	WalletID   int         `json:"wallet_id,omitempty"`
}

type WalletData struct {
	CreateAt       string          `json:"created_at,omitempty"`
	UpdatedAt      string          `json:"updated_at,omitempty"`
	Hold           int             `json:"hold,omitempty"`
	Balance        int             `json:"balance,omitempty"`
	Identification *Identification `json:"identification,omitempty"`
}

type Identification struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func MarshalResponse(w http.ResponseWriter, status int, response *JSONMessage) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(&response)

	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		message, _ := json.Marshal(&JSONMessage{
			Code:    REQUEST_ERROR_CODE,
			Message: err.Error(),
		})
		w.Write(message)
		return
	}

	w.WriteHeader(status)
	w.Write(data)
}
