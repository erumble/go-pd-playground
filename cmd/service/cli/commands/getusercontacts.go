package commands

import (
	"encoding/json"
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/erumble/go-api-boilerplate/cmd/service/cli"
	"github.com/erumble/go-api-boilerplate/pkg/logger"
)

type getUserCmd struct {
	LogLevel string `short:"l" long:"log-level" env:"LOG_LEVEL" default:"debug" required:"false" description:"The log level to use" choice:"debug" choice:"info" choice:"error"`
	APIKey   string `long:"pd-api-key" env:"PD_API_KEY" required:"true" description:"API Key for PagerDuty"`
}

// type getUserResp struct {
// 	Name string `json:"Name"`
// 	ID   string `json:"ID"`
// 	pagerduty.
// }

func init() {
	var cmd getUserCmd
	if _, err := cli.AddCommand(
		&cmd,
		"getuser",
		"Get a user's info in PagerDuty",
		"Get a user's info in PagerDuty(long description)",
	); err != nil {
		// yep, panic. If this fails something is wrong with either the startCmd struct, or the startCmd.Execute function
		panic(err)
	}
}

func (cmd getUserCmd) Execute(args []string) error {
	// Set up the logger
	log := logger.NewFatalLogger(cmd.LogLevel)
	defer log.Sync()
	log.Debug("DEBUG logging enabled")

	pdClient := pagerduty.NewClient(cmd.APIKey)

	resp, err := pdClient.ListUserContactMethods(args[0])
	if err != nil {
		log.Fatal(err)
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string([]byte(respBytes)))
	return nil
}
