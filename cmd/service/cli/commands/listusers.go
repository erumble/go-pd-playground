package commands

import (
	"encoding/json"
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/erumble/go-pd-playground/cmd/service/cli"
	"github.com/erumble/go-pd-playground/pkg/logger"
)

type listUserCmd struct {
	LogLevel string `short:"l" long:"log-level" env:"LOG_LEVEL" default:"debug" required:"false" description:"The log level to use" choice:"debug" choice:"info" choice:"error"`
	APIKey   string `long:"pd-api-key" env:"PD_API_KEY" required:"true" description:"API Key for PagerDuty"`
}

type entry struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
}

func init() {
	var cmd listUserCmd
	if _, err := cli.AddCommand(
		&cmd,
		"listusers",
		"List users in PagerDuty",
		"List users in PagerDuty(long description)",
	); err != nil {
		// yep, panic. If this fails something is wrong with either the startCmd struct, or the startCmd.Execute function
		panic(err)
	}
}

func (cmd listUserCmd) Execute(_ []string) error {
	// Set up the logger
	log := logger.NewFatalLogger(cmd.LogLevel)
	defer log.Sync()
	log.Debug("DEBUG logging enabled")

	pdClient := pagerduty.NewClient(cmd.APIKey)

	resp, err := pdClient.ListUsers(pagerduty.ListUsersOptions{})
	if err != nil {
		log.Fatal(err)
	}

	entries := []entry{}
	for _, r := range resp.Users {
		log.Debugf("user: %s", r.ID)
		e := entry{
			Name: r.Name,
			ID:   r.ID,
		}

		log.Debugf("e: %v", e)

		entries = append(entries, e)
	}

	respBytes, err := json.Marshal(entries)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string([]byte(respBytes)))
	return nil
}
