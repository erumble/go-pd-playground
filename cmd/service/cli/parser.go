package cli

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

// GlobalOpts defines all global command line flags.
var GlobalOpts struct{}

var optParser = flags.NewParser(&GlobalOpts, flags.HelpFlag|flags.PassDoubleDash)

// Execute wraps the flags.Parser.Parse() method and eats the error if it is the help message.
func Execute() error {
	if _, err := optParser.Parse(); err != nil {
		// If the error is the help message, print it and eat it; it's a valid use case.
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Println(err)
			return nil
		}

		// If the error isn't the help message, bubble it up.
		return err
	}
	return nil
}

// AddCommand wraps the flags.Parser.AddCommand() method.
func AddCommand(cmd flags.Commander, name, shortDescription, longDescription string) (*flags.Command, error) {
	return optParser.AddCommand(name, shortDescription, longDescription, cmd)
}
