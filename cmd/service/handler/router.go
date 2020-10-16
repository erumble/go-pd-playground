package handler

import (
	"net/http"

	"github.com/PagerDuty/go-pagerduty"

	"github.com/erumble/go-pd-playground/pkg/logger"
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

// TODO: make interface for pagerduty.Client

// New registers the routes and middleware for the server and returns an http handler
func New(pdClient *pagerduty.Client, log logger.LeveledLogger) Router {
	r := mux.NewRouter()
	r.HandleFunc("/listusers", listUsersHandler(pdClient, log))
	r.PathPrefix("/").Handler(echoHandler(log))
	return &router{r}
}

func (r *router) WithMiddleware(middleware ...mux.MiddlewareFunc) {
	r.Use(middleware...)
}
