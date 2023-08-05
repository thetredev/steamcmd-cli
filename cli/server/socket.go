package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thetredev/steamcmd-cli/shared"
)

func verifyConfiguration() {
	if len(shared.SocketConfig.SocketIp) == 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_IP not set")
	}

	if shared.SocketConfig.SocketPort <= 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_PORT not set")
	}
}

func loadServerCertificates() *tls.Config {
	cert, err := tls.LoadX509KeyPair(ServerCertificates.CertificatePath, ServerCertificates.CertificateKeyPath)

	if err != nil {
		log.Fatalf("Could not load server key pair: %s", err)
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
}

func dialSocket(config *tls.Config) *tls.Conn {
	socket, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", shared.SocketConfig.SocketIp, shared.SocketConfig.SocketPort), config)

	if err != nil {
		log.Fatalf("Could not establish connection: %s", err)
	}

	return socket
}

func sendMessageToSocket(socket *tls.Conn, message string, args ...string) {
	command := []string{message}
	command = append(command, args...)

	fmt.Fprintln(socket, strings.Join(command, " "))
}

func SendMessage(message string, args ...string) {
	verifyConfiguration()

	config := loadServerCertificates()
	socket := dialSocket(config)

	sendMessageToSocket(socket, message, args...)
	reader := bufio.NewReader(socket)

	download := false

	for {
		count, buffer, err := shared.ReadBuffer(reader, 256)

		if err != nil {
			fmt.Printf("ERR SOCKET READ: %s\n", err.Error())
			break
		}

		if count == 0 {
			break
		}

		message := strings.Split(string(buffer[:count]), "\n")[0]

		if strings.HasPrefix(message, shared.MESSAGE_SOCKET_END) {
			break
		} else if strings.HasPrefix(message, shared.MESSAGE_SERVER_FILE_TRANSFER_START) {
			download = true

			fileInfo := strings.Split(message, ";")[1:]
			filePath := filepath.Join(args[1], fileInfo[0])
			fileSize, err := strconv.ParseInt(fileInfo[1], 10, 64)

			if err != nil {
				fmt.Printf("ERR FILE SIZE CONV: %s\n", err.Error())
				break
			}

			shared.ReceiveFile(reader, filePath, fileSize)
		} else if download {
			download = false
		} else {
			fmt.Println(message)
		}
	}

	socket.Close()
}
