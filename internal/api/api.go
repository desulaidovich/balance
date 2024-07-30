package api

import (
	"net/http"

	"github.com/desulaidovich/balance/internal/models"
	"github.com/desulaidovich/balance/internal/services"
	"github.com/desulaidovich/balance/internal/utils"
	"github.com/desulaidovich/balance/pkg/messaging"
	"github.com/desulaidovich/balance/pkg/slogger"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
)

type HttpApi struct {
	mux     *http.ServeMux
	db      *sqlx.DB
	service *services.Service
	nc      *messaging.NatsConnection
	slogger *slogger.Logger
}

func New(mux *http.ServeMux, db *sqlx.DB, nc *nats.Conn, slogger *slogger.Logger) *HttpApi {
	s := services.New(db)
	n := messaging.NewNatsConnection(nc)

	return &HttpApi{
		mux:     mux,
		db:      db,
		service: s,
		nc:      n,
		slogger: slogger,
	}
}

func (h *HttpApi) Create(w http.ResponseWriter, r *http.Request) {
	money, err := utils.GetParamsByName(r, "money")
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	level, err := utils.GetParamsByName(r, "level")
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	wallet := &models.Wallet{
		Balance:             money,
		Hold:                0,
		IdentificationLevel: level,
	}

	limit, err := h.service.GetLimitByID(level)

	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err := wallet.LimitLawCheck(limit); err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err = h.service.CreateWallet(wallet); err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	message := utils.JSONMessage{
		Code: http.StatusOK,
		Data: &utils.Data{
			WalletID: wallet.ID,
			WalletData: &utils.WalletData{
				Balance: wallet.Balance,
				Identification: &utils.Identification{
					ID:   wallet.IdentificationLevel,
					Name: limit.IdentificationLevel,
				},
			},
		},
	}

	if err = h.nc.SendJSON("created", message); err != nil {
		h.slogger.Error(err.Error())
	}

	utils.MarshalResponse(w, http.StatusOK, &utils.JSONMessage{
		Code:    utils.REQUEST_NO_ERROR_CODE,
		Message: "ok",
	})
}

func (h *HttpApi) Hold(w http.ResponseWriter, r *http.Request) {
	walletID, err := utils.GetParamsByName(r, "wallet_id")
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	money, err := utils.GetParamsByName(r, "money")
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err = wallet.HoldBalance(money); err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err = h.service.UpdateWallet(wallet); err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	message := utils.JSONMessage{
		Code: utils.REQUEST_NO_ERROR_CODE,
		Data: &utils.Data{
			WalletID: wallet.ID,
			WalletData: &utils.WalletData{
				Balance: wallet.Balance,
				Hold:    money,
			},
		},
	}

	if err = h.nc.SendJSON("holded", message); err != nil {
		h.slogger.Error(err.Error())
	}

	utils.MarshalResponse(w, http.StatusOK, &utils.JSONMessage{
		Code:    utils.REQUEST_NO_ERROR_CODE,
		Message: "ok",
	})
}

func (h *HttpApi) Dishold(w http.ResponseWriter, r *http.Request) {
	walletID, err := utils.GetParamsByName(r, "wallet_id")
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	money, err := utils.GetParamsByName(r, "money")
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err = wallet.DisholdBalance(money); err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err = h.service.UpdateWallet(wallet); err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	message := utils.JSONMessage{
		Code: utils.REQUEST_NO_ERROR_CODE,
		Data: &utils.Data{
			WalletID: walletID,
			WalletData: &utils.WalletData{
				Hold: wallet.Hold,
			},
		},
	}

	if err = h.nc.SendJSON("disholded", message); err != nil {
		h.slogger.Error(err.Error())
	}

	utils.MarshalResponse(w, http.StatusOK, &utils.JSONMessage{
		Code:    utils.REQUEST_NO_ERROR_CODE,
		Message: "ok",
	})
}

func (h *HttpApi) Edit(w http.ResponseWriter, r *http.Request) {
	walletID, err := utils.GetParamsByName(r, "wallet_id")
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	money, err := utils.GetParamsByName(r, "money")
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	typeID, err := utils.GetParamsByName(r, "type_id")
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	limit, err := h.service.GetLimitByID(wallet.IdentificationLevel)
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}
	if err := wallet.EditWithType(limit, typeID, money); err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	message := utils.JSONMessage{
		Code: utils.REQUEST_NO_ERROR_CODE,
		Data: &utils.Data{
			WalletID: wallet.ID,
			WalletData: &utils.WalletData{
				Balance: wallet.Balance,
				Hold:    wallet.Hold,
			},
		},
	}

	if err = h.nc.SendJSON("edited", message); err != nil {
		h.slogger.Error(err.Error())
	}

	utils.MarshalResponse(w, http.StatusOK, &utils.JSONMessage{
		Code:    utils.REQUEST_NO_ERROR_CODE,
		Message: "ok",
	})
}

func (h *HttpApi) Get(w http.ResponseWriter, r *http.Request) {
	walletID, err := utils.GetParamsByName(r, "wallet_id")
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		message := utils.JSONMessage{
			Code:    utils.REQUEST_ERROR_CODE,
			Message: err.Error(),
		}
		utils.MarshalResponse(w, http.StatusBadRequest, &message)
		return
	}

	createdAt, updatedAt := wallet.GetDates()

	message := utils.JSONMessage{
		Code: utils.REQUEST_NO_ERROR_CODE,
		Data: &utils.Data{
			WalletID: wallet.ID,
			WalletData: &utils.WalletData{
				CreateAt:  createdAt,
				UpdatedAt: updatedAt,
				Balance:   wallet.Balance,
				Identification: &utils.Identification{
					ID: wallet.IdentificationLevel,
				},
			},
		},
	}
	if err = h.nc.SendJSON("got", message); err != nil {
		h.slogger.Error(err.Error())
	}

	utils.MarshalResponse(w, http.StatusOK, &utils.JSONMessage{
		Code:    utils.REQUEST_NO_ERROR_CODE,
		Message: "ok",
	})
}
