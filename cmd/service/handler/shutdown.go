package handler

import (
	"context"
	"net/http"

	"github.com/erumble/go-api-boilerplate/pkg/logger"
)

func shutdownHandler(cancel context.CancelFunc, log logger.LeveledLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shutdown(w, r, cancel, log)
	}
}

func shutdown(w http.ResponseWriter, r *http.Request, cancel context.CancelFunc, log logger.LeveledLogger) {
	log.Debug("Shutting down server")
	w.Write([]byte("OK"))
	cancel()
}
