package balance

import (
	"net/http"
	"strings"

	"github.com/desulaidovich/balance/config"
	"github.com/desulaidovich/balance/internal/api"
	"github.com/desulaidovich/balance/pkg/db"
	"github.com/desulaidovich/balance/pkg/messaging"
	"github.com/desulaidovich/balance/pkg/slogger"
)

func Run() error {
	logger := slogger.New()

	cfg, err := config.LoadEnvFromFile()
	if err != nil {
		return err
	}

	nats, err := messaging.Connect()
	if err != nil {
		return err
	}
	defer nats.Close()

	postgres, err := db.New(cfg)
	if err != nil {
		return err
	}
	defer postgres.Close()

	mux := http.NewServeMux()
	httpApi := api.New(mux, postgres, nats, logger)
	mux.HandleFunc("POST /wallet/create", httpApi.Create)
	mux.HandleFunc("POST /wallet/hold", httpApi.Hold)
	mux.HandleFunc("POST /wallet/dishold", httpApi.Dishold)
	mux.HandleFunc("POST /wallet/edit", httpApi.Edit)
	mux.HandleFunc("GET /wallet/get", httpApi.Get)

	server := new(http.Server)
	server.Addr = ":" + cfg.Port
	server.Handler = logger.Init(mux)

	// Хз)
	logger.Info(strings.ToUpper(cfg.Name) + " started on http://localhost:" + cfg.Port)

	if err = server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
