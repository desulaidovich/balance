package balance

import (
	"log"
	"net/http"

	"github.com/desulaidovich/balance/config"
	"github.com/desulaidovich/balance/internal/api"
	"github.com/desulaidovich/balance/pkg/db"
)

func Run() {
	cfg := config.Load()

	postgres, err := db.NewPostgres(cfg)

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	httpApi := api.New(mux, postgres)

	mux.HandleFunc("POST /wallet/create", httpApi.Create)
	mux.HandleFunc("POST /wallet/hold", httpApi.Hold)
	mux.HandleFunc("POST /wallet/dishold", httpApi.Dishold)
	mux.HandleFunc("POST /wallet/edit", httpApi.Edit)
	// mux.HandleFunc("GET /wallet/get", httpApi.Get)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
