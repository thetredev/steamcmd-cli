package daemon

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/thetredev/steamcmd-cli/shared"
)

var socket *net.UDPConn

func SendSocketResponseMessage(receiver *net.UDPAddr, message string) {
	_, err := socket.WriteToUDP([]byte(fmt.Sprintf("%s\n", message)), receiver)

	if err != nil {
		log.Fatal(err)
	}
}

func StartSocket() {
	if shared.Config.SocketPort <= 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_PORT not set")
	}

	addr := net.UDPAddr{
		Port: shared.Config.SocketPort,
		IP:   net.ParseIP("0.0.0.0"),
	}

	var err error
	socket, err = net.ListenUDP("udp", &addr)

	if err != nil {
		log.Fatal(err)
	}

	server := NewServer()
	server.Logger.Printf("Daemon socket port: %d\n", socket.LocalAddr().(*net.UDPAddr).Port)
	server.Logger.Println("Listening for incoming requests...")

	for {
		buffer := make([]byte, 256)
		_, receiver, err := socket.ReadFromUDP(buffer)

		if err != nil {
			log.Fatal(err)
		}

		message := strings.Split(string(buffer), "\n")[0]

		switch message {
		case "logs":
			server.SendLogs(receiver)
		case "start":
			if err := server.Start(receiver); err != nil {
				log.Fatal(err)
			}
		case "stop":
			server.Stop(receiver)
		case "update":
			if err := server.Update(receiver); err != nil {
				log.Fatal(err)
			}
		default:
			if !handleSpecialMessage(server, receiver, message) {
				SendSocketResponseMessage(receiver, fmt.Sprintf("Invalid command: %s; ignoring...\n", message))
			}
		}

		SendSocketResponseMessage(receiver, shared.SocketEndMessage)
	}
}

func handleSpecialMessage(serverInstance *Server, receiver *net.UDPAddr, message string) bool {
	if strings.HasPrefix(message, "command") {
		command := strings.Join(strings.Split(message, " ")[1:], " ")
		serverInstance.DispatchConsoleCommand(receiver, command)
		return true
	}

	return false
}
