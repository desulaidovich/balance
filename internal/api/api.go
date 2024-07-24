package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/desulaidovich/balance/internal/jsonutil"
	"github.com/desulaidovich/balance/internal/models"
	"github.com/desulaidovich/balance/internal/services"
	"github.com/jmoiron/sqlx"
)

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
	moneyParam := r.URL.Query().Get("money")
	money, err := strconv.Atoi(moneyParam)

	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: `Параметр "moeny" должен быть числом`,
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	levelParam := r.URL.Query().Get("level")
	level, err := strconv.Atoi(levelParam)

	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: `Параметр "level" должен быть числом`,
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	wallet := &models.Wallet{
		Balance:             money,
		Hold:                0,
		IdentificationLevel: level,
	}

	limit, err := h.service.GetLimitByID(levelParam)

	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	if !wallet.LimitLawCheck(limit) {
		limitMaxValue := limit.BalanceMax
		limitMinValue := limit.BalanceMin
		limitName := limit.IdentificationLevel

		text := fmt.Sprintf("В тарифе %s может быть от %d до %d руб.",
			limitName, limitMinValue, limitMaxValue)

		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: text,
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	if err = h.service.CreateWallet(wallet); err != nil {
		message := jsonutil.JsonMessage{
			Code: http.StatusBadRequest,
			// Лень придумывать
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	message := jsonutil.JsonMessage{
		Code:    http.StatusOK,
		Message: "OK",
		Node: map[string]any{
			"id":       wallet.ID,
			"balance":  wallet.Balance,
			"currency": "rub",
			"identification_level": map[string]any{
				"id":   wallet.IdentificationLevel,
				"name": limit.IdentificationLevel,
			},
		},
	}
	jsonutil.MarshalResponse(w, http.StatusOK, "success", &message)
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
