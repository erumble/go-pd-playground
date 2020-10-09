package handler

import (
	"net/http"

	"github.com/erumble/go-api-boilerplate/pkg/logger"
)

func echoHandler(log logger.LeveledLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		echo(w, r, log)
	}
}

func echo(w http.ResponseWriter, r *http.Request, log logger.LeveledLogger) {
	log.Debugf("Echoing back request made to %s to client (%s)", r.URL.Path, r.RemoteAddr)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Allow pre-flight headers
	w.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")

	r.Write(w)
}
