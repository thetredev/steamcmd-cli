package daemon

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/thetredev/steamcmd-cli/shared"
)

const tcpCongestionPreventionDelay = 1 * time.Millisecond

type Socket struct {
	Listener   net.Listener
	Connection net.Conn
}

func NewSocket(ip string, port int, caPath string, certPath string, keyPath string) *Socket {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)

	if err != nil {
		log.Fatalf("Could not load daemon key pair: %s", err)
	}

	certpool := x509.NewCertPool()
	pem, err := os.ReadFile(caPath)

	if err != nil {
		log.Fatalf("Failed to read certificate authority: %v", err)
	}

	if !certpool.AppendCertsFromPEM(pem) {
		log.Fatalf("Could not parse certificate authority")
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certpool,
		Rand:         rand.Reader,
	}

	addr := fmt.Sprintf("%s:%d", ip, port)

	listener, err := tls.Listen("tcp", addr, &config)

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

func StartSocket(caPath string, certPath string, keyPath string) {
	if shared.Config.SocketPort <= 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_PORT not set")
	}

	socket := NewSocket("0.0.0.0", shared.Config.SocketPort, caPath, certPath, keyPath)

	server := NewServer()
	server.Logger.Printf("Listening for incoming requests on port %d/TCP...\n", shared.Config.SocketPort)

	for {
		var err error
		socket.Connection, err = socket.Listener.Accept()

		conn, ok := socket.Connection.(*tls.Conn)

		if !ok {
			server.Logger.Println("Could not establish connection.")
			socket.Connection.Close()
			continue
		}

		err = conn.Handshake()

		if err != nil {
			server.Logger.Printf("Handshake failed: %s\n", err.Error())
			socket.Connection.Close()
			continue
		}

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
