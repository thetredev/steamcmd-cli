package daemon

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/thetredev/steamcmd-cli/shared"
)

const tcpCongestionPreventionDelay = 1 * time.Millisecond

type Socket struct {
	Listener   *net.TCPListener
	Connection *net.TCPConn
}

func NewSocket(ip string, port int) *Socket {
	addr := net.TCPAddr{
		Port: shared.Config.SocketPort,
		IP:   net.ParseIP(ip),
	}

	listener, err := net.ListenTCP("tcp", &addr)

	if err != nil {
		log.Fatal(err)
	}

	return &Socket{
		Listener:   listener,
		Connection: nil,
	}
}

func (socket *Socket) SendMessage(message string) {
	_, err := socket.Connection.Write([]byte(fmt.Sprintln(message)))

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
	server.Logger.Printf("Listening for incoming requests on port %d/TCP...\n", shared.Config.SocketPort)

	for {
		var err error
		socket.Connection, err = socket.Listener.AcceptTCP()

		buffer := make([]byte, 256)
		_, err = socket.Connection.Read(buffer)

		if err != nil {
			server.Logger.Printf("Ignoring socket error: %s\n", err.Error())
			socket.Connection.Close()
			continue
		}

		message := strings.Split(string(buffer), "\n")[0]

		switch message {
		case shared.ServerLogsMessage:
			server.SendLogs(socket)
		case shared.ServerStartMessage:
			if err := server.Start(socket); err != nil {
				socket.SendMessage(err.Error())
			} else {
				socket.SendMessage("Game server started. You can now view its logs.")
			}
		case shared.ServerStopMessage:
			server.Stop(socket)

			time.Sleep(tcpCongestionPreventionDelay)
			socket.SendMessage("Game server stopped.")
		case shared.ServerUpdateMessage:
			for {
				success, err := server.Update(socket)

				if err != nil {
					log.Fatal(err)
				}

				if success {
					break
				}
			}
		default:
			if !handleSpecialMessage(server, socket, message) {
				socket.SendMessage(fmt.Sprintf("Invalid command: %s; ignoring...\n", message))
			}
		}

		time.Sleep(tcpCongestionPreventionDelay)
		socket.SendMessage(shared.SocketEndMessage)

		socket.Connection.Close()
	}
}

func handleSpecialMessage(serverInstance *Server, socket *Socket, message string) bool {
	if strings.HasPrefix(message, shared.ServerConsoleCommandMessage) {
		command := strings.Join(strings.Split(message, " ")[1:], " ")
		serverInstance.DispatchConsoleCommand(socket, command)
		return true
	}

	return false
}
