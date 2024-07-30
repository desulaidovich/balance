package balance

import (
	"net/http"

	"github.com/desulaidovich/balance/config"
	"github.com/desulaidovich/balance/internal/api"
	"github.com/desulaidovich/balance/pkg/db"
	"github.com/desulaidovich/balance/pkg/messaging"
	"github.com/desulaidovich/balance/pkg/slogger"
)

func Run() error {
	logger := slogger.New()
	cfg, err := config.Load()

	if err != nil {
		return err
	}

	natsConn, err := messaging.Connect()

	if err != nil {
		return err
	}
	defer natsConn.Close()

	postgres, err := db.NewPostgres(cfg)

	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	httpApi := api.New(mux, postgres, natsConn, logger)

	mux.HandleFunc("POST /wallet/create", httpApi.Create)
	mux.HandleFunc("POST /wallet/hold", httpApi.Hold)
	mux.HandleFunc("POST /wallet/dishold", httpApi.Dishold)
	mux.HandleFunc("POST /wallet/edit", httpApi.Edit)
	mux.HandleFunc("GET /wallet/get", httpApi.Get)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: logger.Init(mux),
	}

	if err = server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
