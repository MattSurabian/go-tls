package main

import (
	"github.com/mattsurabian/go-tls/server/command"
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
		"config": func() (cli.Command, error) {
			return &command.GenConfigCommand{
				UI: ui,
			}, nil
		},
		"start": func() (cli.Command, error) {
			return &command.StartCommand{
				UI: ui,
			}, nil
		},
	}
}
