package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/erumble/go-api-boilerplate/pkg/logger"
)

type PDUserClient interface {
	ListUsers(o pagerduty.ListUsersOptions) (*pagerduty.ListUsersResponse, error)
}

func listUsersHandler(pdClient PDUserClient, log logger.LeveledLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		listUsers(w, r, pdClient, log)
	}
}

func listUsers(w http.ResponseWriter, r *http.Request, pdClient PDUserClient, log logger.LeveledLogger) {
	resp, err := pdClient.ListUsers(pagerduty.ListUsersOptions{})
	if err != nil {
		log.Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
