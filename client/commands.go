package main

import (
	"github.com/mattsurabian/go-tls/client/command"
	"github.com/mitchellh/cli"
	"os"
)

// Commands is the mapping of all available commands
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	Commands = map[string]cli.CommandFactory{
		"send": func() (cli.Command, error) {
			return &command.SendCommand{
				UI: ui,
			}, nil
		},
	}
}
