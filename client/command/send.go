package command

import (
	"github.com/mattsurabian/go-tls/shared/tlsUtils"
	"github.com/mitchellh/cli"

	"log"
	"strings"
)

// SendCommand sends data to the server
type SendCommand struct {
	UI cli.Ui
}

// Long-form help
func (c *SendCommand) Help() string {
	help := `
Usage: [flags] send [text]
`
	return strings.TrimSpace(help)
}

func (c *SendCommand) Synopsis() string {
	return "Send data to the server"
}

// Run the actual command
func (c *SendCommand) Run(args []string) int {

	if len(args) < 1 {
		log.Println("Error: Missing arguments, run -h for more info")
		return BAD_REQUEST
	}

	textToSend := []byte(args[0])

	conn, err := tlsUtils.GetClientTLSConnection()
	if err != nil {
		c.UI.Error(err.Error())
	} else {
		conn.Write(textToSend)
	}

	return OK
}
