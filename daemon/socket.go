package daemon

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/thetredev/steamcmd-cli/shared"
)

type Socket struct {
	Connection *net.UDPConn
}

func NewSocket(ip string, port int) *Socket {
	addr := net.UDPAddr{
		Port: shared.Config.SocketPort,
		IP:   net.ParseIP(ip),
	}

	connection, err := net.ListenUDP("udp", &addr)

	if err != nil {
		log.Fatal(err)
	}

	return &Socket{
		Connection: connection,
	}
}

func (socket *Socket) SendMessage(receiver *net.UDPAddr, message string) {
	_, err := socket.Connection.WriteToUDP([]byte(fmt.Sprintf("%s\n", message)), receiver)

	if err != nil {
		log.Fatal(err)
	}
}

func StartSocket() {
	if shared.Config.SocketPort <= 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_PORT not set")
	}

	socket := NewSocket("0.0.0.0", shared.Config.SocketPort)

	server := NewServer()
	server.Logger.Printf("Daemon socket port: %d\n", socket.Connection.LocalAddr().(*net.UDPAddr).Port)
	server.Logger.Println("Listening for incoming requests...")

	for {
		buffer := make([]byte, 256)
		_, receiver, err := socket.Connection.ReadFromUDP(buffer)

		if err != nil {
			log.Fatal(err)
		}

		message := strings.Split(string(buffer), "\n")[0]

		switch message {
		case shared.ServerLogsMessage:
			server.SendLogs(socket, receiver)
		case shared.ServerStartMessage:
			if err := server.Start(socket, receiver); err != nil {
				log.Fatal(err)
			}
		case shared.ServerStopMessage:
			server.Stop(socket, receiver)
		case shared.ServerUpdateMessage:
			if err := server.Update(socket, receiver); err != nil {
				log.Fatal(err)
			}
		default:
			if !handleSpecialMessage(server, socket, receiver, message) {
				socket.SendMessage(receiver, fmt.Sprintf("Invalid command: %s; ignoring...\n", message))
			}
		}

		socket.SendMessage(receiver, shared.SocketEndMessage)
	}
}

func handleSpecialMessage(serverInstance *Server, socket *Socket, receiver *net.UDPAddr, message string) bool {
	if strings.HasPrefix(message, shared.ServerCommandMessage) {
		command := strings.Join(strings.Split(message, " ")[1:], " ")
		serverInstance.DispatchConsoleCommand(socket, receiver, command)
		return true
	}

	return false
}
