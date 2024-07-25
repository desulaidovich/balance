package api

import (
	"net/http"

	"github.com/desulaidovich/balance/internal/jsonutil"
	"github.com/desulaidovich/balance/internal/models"
	"github.com/desulaidovich/balance/internal/requtil"
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
	money, err := requtil.GetParamsByName(r, "money")
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: `параметр "moeny" должен быть числом`,
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	level, err := requtil.GetParamsByName(r, "level")
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: `параметр "level" должен быть числом`,
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	wallet := &models.Wallet{
		Balance:             money,
		Hold:                0,
		IdentificationLevel: level,
	}

	limit, err := h.service.GetLimitByID(level)

	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err := wallet.LimitLawCheck(limit); err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err = h.service.CreateWallet(wallet); err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	message := jsonutil.JsonMessage{
		Code: http.StatusOK,
		Data: map[string]any{
			"wallet_id": wallet.ID,
			"wallet_data": map[string]any{
				"balance": wallet.Balance,
				"identification": map[string]any{
					"id":   wallet.IdentificationLevel,
					"name": limit.IdentificationLevel,
				},
			},
		},
	}
	jsonutil.MarshalResponse(w, http.StatusOK, &message)
}

func (h *HttpApi) Hold(w http.ResponseWriter, r *http.Request) {
	walletID, err := requtil.GetParamsByName(r, "wallet_id")
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	money, err := requtil.GetParamsByName(r, "money")
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err = wallet.HoldBalance(money); err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err = h.service.UpdateWallet(wallet); err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	message := jsonutil.JsonMessage{
		Code: http.StatusOK,
		Data: map[string]any{
			"hold": wallet.Hold,
		},
	}
	jsonutil.MarshalResponse(w, http.StatusOK, &message)
}

func (h *HttpApi) Dishold(w http.ResponseWriter, r *http.Request) {
	walletID, err := requtil.GetParamsByName(r, "wallet_id")
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	money, err := requtil.GetParamsByName(r, "money")
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err = wallet.DisholdBalance(money); err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err = h.service.UpdateWallet(wallet); err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	message := jsonutil.JsonMessage{
		Code: http.StatusOK,
		Data: map[string]any{
			"hold": wallet.Hold,
		},
	}
	jsonutil.MarshalResponse(w, http.StatusOK, &message)
}

func (h *HttpApi) Edit(w http.ResponseWriter, r *http.Request) {
	walletID, err := requtil.GetParamsByName(r, "wallet_id")
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	money, err := requtil.GetParamsByName(r, "money")
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	typeID, err := requtil.GetParamsByName(r, "type_id")
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	limit, err := h.service.GetLimitByID(wallet.IdentificationLevel)
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err := wallet.EditWithType(limit, typeID, money); err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	message := jsonutil.JsonMessage{
		Code: http.StatusOK,
		Data: map[string]any{
			"balance": wallet.Balance,
			"hold":    wallet.Hold,
		},
	}
	jsonutil.MarshalResponse(w, http.StatusOK, &message)
}

func (h *HttpApi) Get(w http.ResponseWriter, r *http.Request) {
	walletID, err := requtil.GetParamsByName(r, "wallet_id")
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		message := jsonutil.JsonMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		jsonutil.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	createdAt, updatedAt := wallet.GetDates()

	message := jsonutil.JsonMessage{
		Code: http.StatusOK,
		Data: map[string]any{
			"wallet_id": wallet.ID,
			"wallet_data": map[string]any{
				"banalce": wallet.Balance,
				"hold":    wallet.Hold,
				"date": map[string]any{
					"created": createdAt,
					"updated": updatedAt,
				},
			},
		},
	}
	jsonutil.MarshalResponse(w, http.StatusOK, &message)
}
