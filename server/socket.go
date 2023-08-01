package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
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

	cert, err := tls.LoadX509KeyPair("/certs/server/cert.pem", "/certs/server/cert.key")

	if err != nil {
		log.Fatalf("Could not load server key pair: %s", err)
	}

	config := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	socket, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", shared.Config.SocketIp, shared.Config.SocketPort), &config)

	if err != nil {
		log.Fatalf("Could not establish connection: %s", err)
	}

	command := []string{message}
	command = append(command, args...)

	fmt.Fprintln(socket, strings.Join(command, " "))
	reader := bufio.NewReader(socket)

	for {
		buffer := make([]byte, 2048)
		_, err := reader.Read(buffer)

		if err != nil {
			break
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
