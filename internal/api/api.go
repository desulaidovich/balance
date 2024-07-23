package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/desulaidovich/balance/internal/services"
	"github.com/jmoiron/sqlx"
)

type JsonMessage struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

type HttpApi struct {
	mux     *http.ServeMux
	db      *sqlx.DB
	service *services.Service
}

func New(mux *http.ServeMux, db *sqlx.DB) *HttpApi {
	s := services.New(db)

	return &HttpApi{
		mux:     mux,
		db:      db,
		service: s,
	}
}

func (h *HttpApi) Create(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("money")
	money, err := strconv.Atoi(param)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		jsonError := &JsonMessage{
			Code:    http.StatusBadRequest,
			Message: `Параметр "money" должен быть числом`,
		}

		request, _ := json.Marshal(map[string]JsonMessage{
			"error": *jsonError,
		})

		w.Write(request)
		return
	}

	if money <= 0 {
		w.WriteHeader(http.StatusBadRequest)

		jsonError := &JsonMessage{
			Code:    http.StatusBadRequest,
			Message: `Параметр "money" должен быть больше нуля`,
		}

		request, _ := json.Marshal(map[string]JsonMessage{
			"error": *jsonError,
		})

		w.Write(request)
		return
	}

	h.service.CreateWallet()

	jsonError := &JsonMessage{
		Code:    http.StatusOK,
		Message: "Создан новый кошелек с суммой " + param + " рублей",
	}

	request, _ := json.Marshal(map[string]JsonMessage{
		"error": *jsonError,
	})

	w.Write(request)
}

// func (h *HttpApi) Hold(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("POST host:port/balance/hold?id=id&money=сумма"))
// }

// func (h *HttpApi) Dishold(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("POST host:port/balance/dishold?id=id&money=сумма"))
// }

// func (h *HttpApi) Edit(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("POST host:port/balance/edit?id=id&money=сумма&type=списание/пополнение"))
// }

// func (h *HttpApi) Get(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("GET host:port/balance/get?id=id"))
// }
