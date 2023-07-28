package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/thetredev/steamcmd-cli/shared"
)

var socket net.Conn

func init() {
	if shared.Config.SocketPort <= 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_PORT not set")
	}

	var err error
	socket, err = net.Dial("udp", fmt.Sprintf("127.0.0.1:%d", shared.Config.SocketPort))

	if err != nil {
		log.Fatal(err)
	}
}

func SendMessageToSocket(message string, args ...string) {
	if len(args) == 0 {
		fmt.Fprintf(socket, "%s\n", message)
	} else {
		fmt.Fprintf(socket, "%s %s\n", message, strings.Join(args, " "))
	}

	reader := bufio.NewReader(socket)

	for {
		buffer := make([]byte, 2048)
		_, err := reader.Read(buffer)

		if err != nil {
			stopSocket()
			log.Fatal(err)
		}

		message := string(buffer)

		if strings.HasPrefix(message, shared.SocketEndMessage) {
			break
		} else {
			fmt.Print(message)
		}
	}
}

// TODO: Call this on SIGTERM/SIGQUIT/SIGKILL
func stopSocket() {
	socket.Close()
}
