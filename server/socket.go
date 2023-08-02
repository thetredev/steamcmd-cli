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
	if len(shared.SocketConfig.SocketIp) == 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_IP not set")
	}

	if shared.SocketConfig.SocketPort <= 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_PORT not set")
	}

	cert, err := tls.LoadX509KeyPair(ServerCertificates.CertificatePath, ServerCertificates.CertificateKeyPath)

	if err != nil {
		log.Fatalf("Could not load server key pair: %s", err)
	}

	config := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	socket, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", shared.SocketConfig.SocketIp, shared.SocketConfig.SocketPort), &config)

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

		if strings.HasPrefix(message, shared.MESSAGE_SOCKET_END) {
			break
		} else {
			fmt.Print(message)
		}
	}

	socket.Close()
}
