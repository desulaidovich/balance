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

// /wallet/create?money=INT_VALUE&level=INT_VALUE
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
		Message: "created",
		Node: map[string]any{
			"wallet_id": wallet.ID,
			"balance":   wallet.Balance,
			"currency":  "rub",
			"identification_level": map[string]any{
				"id":   wallet.IdentificationLevel,
				"name": limit.IdentificationLevel,
			},
		},
	}
	jsonutil.MarshalResponse(w, http.StatusOK, "success", &message)
}

// /wallet/hold?wallet_id=INT_VALUE&money=INT_VALUE
func (h *HttpApi) Hold(w http.ResponseWriter, r *http.Request) {
	walletIDParam := r.URL.Query().Get("wallet_id")
	walletID, err := strconv.Atoi(walletIDParam)

	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: `Параметр "wallet_id" должен быть числом`,
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	moneyParam := r.URL.Query().Get("money")
	money, err := strconv.Atoi(moneyParam)

	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: `Параметр "money" должен быть числом`,
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)

	if err != nil {
		message := jsonutil.JsonMessage{
			Code: http.StatusBadRequest,
			// Лень придумывать
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	if err = wallet.HoldBalance(money); err != nil {
		message := jsonutil.JsonMessage{
			Code: http.StatusBadRequest,
			// Лень придумывать
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	if err = h.service.UpdateWallet(wallet); err != nil {
		message := jsonutil.JsonMessage{
			Code: http.StatusBadRequest,
			// Лень придумывать
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
	}
	message := jsonutil.JsonMessage{
		Code:    http.StatusOK,
		Message: "holded",
		Node: map[string]any{
			"wallet_id": wallet.ID,
			"hold":      wallet.Hold,
		},
	}
	jsonutil.MarshalResponse(w, http.StatusOK, "success", &message)
}

// /wallet/dishold?wallet_id=INT_VALUE&money=INT_VALUE
func (h *HttpApi) Dishold(w http.ResponseWriter, r *http.Request) {
	walletIDParam := r.URL.Query().Get("wallet_id")
	walletID, err := strconv.Atoi(walletIDParam)

	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: `Параметр "wallet_id" должен быть числом`,
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	moneyParam := r.URL.Query().Get("money")
	money, err := strconv.Atoi(moneyParam)

	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: `Параметр "money" должен быть числом`,
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)

	if err != nil {
		message := jsonutil.JsonMessage{
			Code: http.StatusBadRequest,
			// Лень придумывать
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	if err = wallet.DisholdBalance(money); err != nil {
		message := jsonutil.JsonMessage{
			Code: http.StatusBadRequest,
			// Лень придумывать
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
		return
	}

	if err = h.service.UpdateWallet(wallet); err != nil {
		message := jsonutil.JsonMessage{
			Code: http.StatusBadRequest,
			// Лень придумывать
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, "error", &message)
	}

	message := jsonutil.JsonMessage{
		Code: http.StatusOK,
		// Лень придумывать
		Message: "disholded",
		Node: map[string]any{
			"wallet_id":    wallet.ID,
			"disholded":    money,
			"current_hold": wallet.Hold,
		},
	}
	jsonutil.MarshalResponse(w, http.StatusOK, "success", &message)
}

// func (h *HttpApi) Edit(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("POST host:port/balance/edit?id=id&money=сумма&type=списание/пополнение"))
// }

// func (h *HttpApi) Get(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("GET host:port/balance/get?id=id"))
// }
