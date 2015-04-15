package command

import (
	"github.com/mattsurabian/go-tls/shared/tlsUtils"
	"github.com/mitchellh/cli"
	"log"
	"net"
	"strings"
)

// StartCommand starts the server application listening on the configured port
type StartCommand struct {
	UI cli.Ui
}

// Long-form help
func (c *StartCommand) Help() string {
	help := `
Usage: start
  This command will start the server process
`
	return strings.TrimSpace(help)
}

func (c *StartCommand) Synopsis() string {
	return "Start the server"
}

// Run the actual command
func (c *StartCommand) Run(args []string) int {
	listener := tlsUtils.GetServerTLSListener()

	for {
		conn, err := listener.Accept()

		if err != nil {
			panic(err)
		}

		log.Println("------------------------------------")
		log.Println("connection open")
		go handleClient(conn)
	}
	return OK
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		// log output for now, eventually we should store this somewhere
		log.Printf("received: %s\n", buf[:n])
	}
	log.Println("connection closed")
	log.Println("------------------------------------")
}
