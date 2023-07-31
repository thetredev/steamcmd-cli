package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/thetredev/steamcmd-cli/shared"
)

func SendMessage(message string, args ...string) {
	if len(shared.Config.SocketIp) == 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_IP not set")
	}

	if shared.Config.SocketPort <= 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_PORT not set")
	}

	socket, err := net.Dial("tcp", fmt.Sprintf("%s:%d", shared.Config.SocketIp, shared.Config.SocketPort))

	if err != nil {
		log.Fatal(err)
	}

	command := []string{message}
	command = append(command, args...)

	fmt.Fprintln(socket, strings.Join(command, " "))
	reader := bufio.NewReader(socket)

	for {
		buffer := make([]byte, 2048)
		_, err := reader.Read(buffer)

		if err != nil {
			socket.Close()
			log.Fatal(err)
		}

		message := string(buffer)

		if strings.HasPrefix(message, shared.SocketEndMessage) {
			break
		} else {
			fmt.Print(message)
		}
	}

	socket.Close()
}
