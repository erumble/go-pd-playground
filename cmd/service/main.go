package main

import (
	"log"

	"github.com/erumble/go-pd-playground/cmd/service/cli"
	_ "github.com/erumble/go-pd-playground/cmd/service/cli/commands"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
