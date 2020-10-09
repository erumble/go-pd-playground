package handler

import (
	"context"
	"net/http"

	"github.com/erumble/go-api-boilerplate/pkg/logger"
	"github.com/gorilla/mux"
)

// Router allows us to pass in middleware
type Router interface {
	http.Handler
	WithMiddleware(middleware ...mux.MiddlewareFunc)
}

type router struct {
	*mux.Router
}

// New registers the routes and middleware for the server and returns an http handler
func New(cancel context.CancelFunc, log logger.LeveledLogger) Router {
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(echoHandler(log))
	return &router{r}
}

func (r *router) WithMiddleware(middleware ...mux.MiddlewareFunc) {
	r.Use(middleware...)
}
