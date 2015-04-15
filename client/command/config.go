package command

import (
	"github.com/mattsurabian/go-tls/shared/cliUtils"
	"github.com/mitchellh/cli"
	"strings"
)

// GenConfigCommand attempts to write out an INI configuration file
type GenConfigCommand struct {
	UI cli.Ui
}

// Long-form help
func (c *GenConfigCommand) Help() string {
	help := `
Usage: config
  This command will prompt the user for configuration values
  all are optional but any provided will be persisted to disk
`
	return strings.TrimSpace(help)
}

func (c *GenConfigCommand) Synopsis() string {
	return "Create or update application configuration"
}

// Run the actual command
func (c *GenConfigCommand) Run(args []string) int {
	cliUtils.GenerateClientConfig()
	return OK
}
