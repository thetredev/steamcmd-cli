package daemon

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/thetredev/steamcmd-cli/shared"
)

type Socket struct {
	Listener   *net.TCPListener
	Connection *net.TCPConn
	Input      chan string
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
		Input:      make(chan string),
	}
}

func (socket *Socket) Delete() {
	socket.Connection.Close()
	socket.Connection = nil
}

func listenForSocketInput(socket *Socket) {
	const delay = 1 * time.Millisecond

	for {
		if socket.Connection == nil {
			break
		}

		select {
		case message := <-socket.Input:
			_, err := socket.Connection.Write([]byte(fmt.Sprintln(message)))

			if err != nil {
				log.Fatal(err)
			}

			// Prevent TCP congestion
			time.Sleep(delay)
		}
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
			log.Fatal(err)
		}

		go listenForSocketInput(socket)
		message := strings.Split(string(buffer), "\n")[0]

		switch message {
		case shared.ServerLogsMessage:
			server.SendLogs(socket)
		case shared.ServerStartMessage:
			if err := server.Start(socket); err != nil {
				socket.Input <- err.Error()
			} else {
				socket.Input <- "Game server started. You can now view its logs."
			}
		case shared.ServerStopMessage:
			server.Stop(socket)
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
				socket.Input <- fmt.Sprintf("Invalid command: %s; ignoring...\n", message)
			}
		}

		socket.Input <- shared.SocketEndMessage
		socket.Delete()
	}
}

func handleSpecialMessage(serverInstance *Server, socket *Socket, message string) bool {
	if strings.HasPrefix(message, shared.ServerCommandMessage) {
		command := strings.Join(strings.Split(message, " ")[1:], " ")
		serverInstance.DispatchConsoleCommand(socket, command)
		return true
	}

	return false
}
