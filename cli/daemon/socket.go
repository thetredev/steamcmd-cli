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

const TCP_CONGESTION_PREVENTION_DELAY = 1 * time.Millisecond

type Socket struct {
	Listener   net.Listener
	Connection net.Conn
}

func NewSocket() *Socket {
	cert, err := tls.LoadX509KeyPair(Config.CertificatePath, Config.CertificateKeyPath)

	if err != nil {
		log.Fatalf("Could not load daemon key pair: %s", err)
	}

	certpool := x509.NewCertPool()
	ca, err := os.ReadFile(Config.CACertificatePath)

	if err != nil {
		log.Fatalf("Failed to read certificate authority: %v", err)
	}

	if !certpool.AppendCertsFromPEM(ca) {
		log.Fatalf("Could not parse certificate authority")
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certpool,
		Rand:         rand.Reader,
	}

	addr := fmt.Sprintf("%s:%d", shared.SocketConfig.SocketIp, shared.SocketConfig.SocketPort)
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

func (socket *Socket) SendAndLogMessage(server *Server, message string) {
	socket.SendMessage(message)
	server.Logger.Println(message)
}

func StartSocket() {
	if shared.SocketConfig.SocketPort <= 0 {
		log.Fatal("STEAMCMD_CLI_SOCKET_PORT not set")
	}

	socket := NewSocket()

	server := NewServer()
	server.Logger.Printf("Listening for incoming requests on port %d/TCP...\n", shared.SocketConfig.SocketPort)

	for {
		var err error
		socket.Connection, err = socket.Listener.Accept()

		if err != nil {
			log.Fatal("could not listen on socket, aborting...")
		}

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
		case shared.MESSAGE_SERVER_LOGS:
			server.SendLogs(socket)

		case shared.MESSAGE_SERVER_START:
			if err := server.Start(socket); err != nil {
				socket.SendMessage(err.Error())
			}

		case shared.MESSAGE_SERVER_STOP:
			server.Stop(socket)

		case shared.MESSAGE_SERVER_UPDATE:
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

		time.Sleep(TCP_CONGESTION_PREVENTION_DELAY)
		socket.SendMessage(shared.MESSAGE_SOCKET_END)

		socket.Connection.Close()
	}
}

func handleSpecialMessage(serverInstance *Server, socket *Socket, message string) bool {
	if strings.HasPrefix(message, shared.MESSAGE_SERVER_CONSOLE_COMMAND) {
		command := strings.Join(strings.Split(message, " ")[1:], " ")
		serverInstance.DispatchConsoleCommand(socket, command)

		return true
	}

	return false
}
