package slogger

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/lmittmann/tint"
)

type Logger struct {
	*slog.Logger
}

func New() *Logger {
	logger := new(Logger)

	logger.Logger = slog.New(tint.NewHandler(os.Stdout, nil))

	return logger
}

func (l *Logger) Init(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info(r.RemoteAddr,
			slog.String(r.Method, r.URL.String()),
		)
		h.ServeHTTP(w, r)
	})
}
