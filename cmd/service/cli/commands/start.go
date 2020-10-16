package commands

import (
	"context"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/erumble/go-pd-playground/cmd/service/cli"
	"github.com/erumble/go-pd-playground/cmd/service/handler"
	"github.com/erumble/go-pd-playground/pkg/logger"
	"github.com/erumble/go-pd-playground/pkg/middleware/httplogger"
	"github.com/erumble/go-pd-playground/pkg/server"
)

type startCmd struct {
	LogLevel string `short:"l" long:"log-level" env:":LOG_LEVEL" default:"debug" required:"false" description:"The log level to use" choice:"debug" choice:"info" choice:"error"`
	Port     int    `short:"p" long:"port"      env:"PORT"       default:"8080" required:"false" description:"The port on which the service listens"`

	APIKey string `long:"pd-api-key" env:"PD_API_KEY" required:"true" description:"API Key for PagerDuty"`
}

func init() {
	var cmd startCmd
	if _, err := cli.AddCommand(
		&cmd,
		"start",
		"Start the service",
		"Start the service (long description)",
	); err != nil {
		// yep, panic. If this fails something is wrong with either the startCmd struct, or the startCmd.Execute function
		panic(err)
	}
}

// Execute implements the Commander interace from the jessevdk/go-flags package
// We don't (currnetly) care about positional arguments, so we use an `_ []string` to ignore them.
func (cmd startCmd) Execute(_ []string) error {
	// Set up the logger
	log := logger.NewLeveledLogger(cmd.LogLevel)
	defer log.Sync()
	log.Debug("DEBUG logging enabled")

	pdClient := pagerduty.NewClient(cmd.APIKey)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := handler.New(pdClient, log)
	h.WithMiddleware(httplogger.HTTPLogger)
	s := server.New(h, cmd.Port, 5*time.Second, log)

	return s.Serve(ctx)
}
