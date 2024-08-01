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
	nats    *messaging.NatsConnection
	slogger *slogger.Logger
}

func New(mux *http.ServeMux, db *sqlx.DB, natsConn *nats.Conn, slogger *slogger.Logger) *HttpApi {
	httpApi := new(HttpApi)
	httpApi.mux = mux
	httpApi.db = db
	httpApi.service = services.New(db)
	httpApi.nats = messaging.New(natsConn)
	httpApi.slogger = slogger

	return httpApi
}

func (h *HttpApi) Create(w http.ResponseWriter, r *http.Request) {
	getMoneyParam := utils.GetIntParam(r, "money")
	money, err := getMoneyParam(w)
	if err != nil {
		h.slogger.Logger.Error(err.Error())
		return
	}

	getLevelParam := utils.GetIntParam(r, "level")
	level, err := getLevelParam(w)
	if err != nil {
		h.slogger.Logger.Error(err.Error())
		return
	}

	wallet := new(models.Wallet)
	wallet.Balance = money
	wallet.Hold = 0
	wallet.IdentificationLevel = level

	limit, err := h.service.GetLimitByID(level)

	if err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
		return
	}

	if err = wallet.LimitLawCheck(limit); err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
		return
	}

	if err = h.service.CreateWallet(wallet); err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
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

	if err = h.nats.SendJSON("created", message); err != nil {
		h.slogger.Error(err.Error())
	}

	msg := new(utils.JSONMessage)
	msg.Code = utils.REQUEST_NO_ERROR_CODE
	msg.Message = "ok"
	if err = msg.RequestByHTTP(w, http.StatusOK); err != nil {
		h.slogger.Logger.Error(err.Error())
	}
}

func (h *HttpApi) Hold(w http.ResponseWriter, r *http.Request) {
	getWalletIDParam := utils.GetIntParam(r, "wallet_id")
	walletID, err := getWalletIDParam(w)
	if err != nil {
		h.slogger.Logger.Error(err.Error())
		return
	}

	getMoneyParam := utils.GetIntParam(r, "money")
	money, err := getMoneyParam(w)
	if err != nil {
		h.slogger.Logger.Error(err.Error())
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
		return
	}

	if err = wallet.HoldBalance(money); err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
		return
	}

	if err = h.service.UpdateWallet(wallet); err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
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

	if err = h.nats.SendJSON("holded", message); err != nil {
		h.slogger.Error(err.Error())
	}

	msg := new(utils.JSONMessage)
	msg.Code = utils.REQUEST_NO_ERROR_CODE
	msg.Message = "ok"
	if err = msg.RequestByHTTP(w, http.StatusOK); err != nil {
		h.slogger.Logger.Error(err.Error())
	}
}

func (h *HttpApi) Dishold(w http.ResponseWriter, r *http.Request) {
	getWalletIDParam := utils.GetIntParam(r, "wallet_id")
	walletID, err := getWalletIDParam(w)
	if err != nil {
		h.slogger.Logger.Error(err.Error())
		return
	}

	getMoneyParam := utils.GetIntParam(r, "money")
	money, err := getMoneyParam(w)
	if err != nil {
		h.slogger.Logger.Error(err.Error())
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
		return
	}

	if err = wallet.DisholdBalance(money); err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
		return
	}

	if err = h.service.UpdateWallet(wallet); err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
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

	if err = h.nats.SendJSON("disholded", message); err != nil {
		h.slogger.Error(err.Error())
	}

	msg := new(utils.JSONMessage)
	msg.Code = utils.REQUEST_NO_ERROR_CODE
	msg.Message = "ok"
	if err = msg.RequestByHTTP(w, http.StatusOK); err != nil {
		h.slogger.Logger.Error(err.Error())
	}
}

func (h *HttpApi) Edit(w http.ResponseWriter, r *http.Request) {
	getWalletIDParam := utils.GetIntParam(r, "wallet_id")
	walletID, err := getWalletIDParam(w)
	if err != nil {
		h.slogger.Logger.Error(err.Error())
		return
	}

	getMoneyParam := utils.GetIntParam(r, "money")
	money, err := getMoneyParam(w)
	if err != nil {
		h.slogger.Logger.Error(err.Error())
		return
	}

	getTypeIDParam := utils.GetIntParam(r, "type_id")
	typeID, err := getTypeIDParam(w)
	if err != nil {
		h.slogger.Logger.Error(err.Error())
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
		return
	}

	limit, err := h.service.GetLimitByID(wallet.IdentificationLevel)
	if err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
		return
	}

	if err = wallet.EditWithType(limit, typeID, money); err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
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

	if err = h.nats.SendJSON("edited", message); err != nil {
		h.slogger.Error(err.Error())
	}

	msg := new(utils.JSONMessage)
	msg.Code = utils.REQUEST_NO_ERROR_CODE
	msg.Message = "ok"
	if err = msg.RequestByHTTP(w, http.StatusOK); err != nil {
		h.slogger.Logger.Error(err.Error())
	}
}

func (h *HttpApi) Get(w http.ResponseWriter, r *http.Request) {
	getTypeIDParam := utils.GetIntParam(r, "wallet_id")
	walletID, err := getTypeIDParam(w)
	if err != nil {
		h.slogger.Logger.Error(err.Error())
		return
	}

	wallet, err := h.service.GetWalletByID(walletID)
	if err != nil {
		msg := new(utils.JSONMessage)
		msg.Code = utils.REQUEST_ERROR_CODE
		msg.Message = err.Error()
		if err = msg.RequestByHTTP(w, http.StatusBadRequest); err != nil {
			h.slogger.Logger.Error(err.Error())
		}
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

	if err = h.nats.SendJSON("got", message); err != nil {
		h.slogger.Error(err.Error())
	}

	msg := new(utils.JSONMessage)
	msg.Code = utils.REQUEST_NO_ERROR_CODE
	msg.Message = "ok"

	if err = msg.RequestByHTTP(w, http.StatusOK); err != nil {
		h.slogger.Logger.Error(err.Error())
	}
}
