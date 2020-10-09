package handler

import (
	"context"
	"net/http"

	"github.com/erumble/go-api-boilerplate/pkg/logger"
	"github.com/gorilla/mux"
)

// New registers the routes and middleware for the server and returns an http handler
func New(cancel context.CancelFunc, log logger.LeveledLogger) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/echo", echoHandler(log))
	r.HandleFunc("/shutdown", shutdownHandler(cancel, log))
	return r
}
