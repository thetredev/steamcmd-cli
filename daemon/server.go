package daemon

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

const serverConsoleInputDelay time.Duration = 250 * time.Millisecond

type Server struct {
	Logger  *log.Logger
	Command *exec.Cmd
	Console *ServerConsole
}

func NewServer() *Server {
	return &Server{
		Logger:  log.Default(),
		Command: nil,
		Console: NewServerConsole(),
	}
}

func (server *Server) Delete() {
	if server.Console != nil {
		server.Console.Delete()
	}

	server.Command = nil
}

func (server *Server) Update(receiver *net.UDPAddr) error {
	server.Logger.Printf("Received request to update the game server from %v\n", receiver)

	if server.IsRunning() {
		return errors.New("Server is running, cannot update. Ignoring...")
	}

	server.Logger.Println("Updating the game server...")

	if len(Config.Application) == 0 {
		return errors.New("STEAMCMD_SH not set")
	}

	if _, err := os.Stat(Config.Application); os.IsNotExist(err) {
		return err
	}

	if Config.ServerAppId <= 0 {
		return errors.New("STEAMCMD_SERVER_APPID not set")
	}

	if len(Config.ServerHome) == 0 {
		return errors.New("STEAMCMD_SERVER_HOME not set")
	}

	updateCommand := exec.Command(
		"bash", "-c",
		fmt.Sprintf(
			"%s +force_install_dir %s +login anonymous %s +app_update %d validate +quit",
			Config.Application,
			Config.ServerHome,
			Config.ServerAppConfig,
			Config.ServerAppId,
		),
	)

	ptyFile, err := pty.Start(updateCommand)

	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(ptyFile)

		for scanner.Scan() {
			SendSocketResponseMessage(receiver, scanner.Text())
		}
	}()

	updateCommand.Wait()
	return nil
}

func (server *Server) Start(receiver *net.UDPAddr) error {
	server.Logger.Printf("Received request to start the game server from %v\n", receiver)

	if server.IsRunning() {
		return errors.New(fmt.Sprintf("Server already running (PID: %d)", server.Command.Process.Pid))

		//server.Logger.Printf("Ignoring: %s\n", message)
		//SendSocketResponseMessage(receiver, message)
	}

	if _, err := os.Stat(Config.ServerHome); os.IsNotExist(err) {
		return errors.New("STEAMCMD_SERVER_HOME is set to a nonexistent path")
	}

	server.Logger.Println("Starting game server...")

	server.Command = exec.Command(
		"bash", "-c",
		"./srcds_linux -console -game cstrike +ip 0.0.0.0 -port 27019 +maxplayers 12 +map de_dust2 -tickrate 128 -threads 3 -nodev",
	)

	server.Command.Env = os.Environ()
	server.Command.Env = append(server.Command.Env,
		fmt.Sprintf("LD_LIBRARY_PATH=.:./bin:%s", os.Getenv("LD_LIBRARY_PATH")),
		"RESTART=no",
	)

	server.Command.Dir = Config.ServerHome

	var err error
	server.Console.Pty, err = pty.Start(server.Command)

	if err != nil {
		return err
	}

	go server.Console.ListenForInput()
	go server.Console.ListenForOutput()

	return nil
}

func (server *Server) Stop(receiver *net.UDPAddr) {
	server.Logger.Printf("Received request to stop the game server from %v\n", receiver)
	server.DispatchConsoleCommand(receiver, "quit")
}

func (server *Server) IsRunning() bool {
	return server.Command != nil && server.Console != nil
}

func (server *Server) IsSRCDS() bool {
	// check whether Config.ServerHome contains the executable 'srcds_linux'
	return false
}

func (server *Server) IsCSGO() bool {
	// isSRCDS + check whether Config.ServerHome contains the folder 'csgo'
	return false
}

func (server *Server) SendLogs(receiver *net.UDPAddr) {
	server.Logger.Printf("Received request to send game server logs to %v\n", receiver)

	if server.IsRunning() {
		server.Logger.Printf("Sending game server logs to %v...\n", receiver)

		bytes := server.Console.SendLogs(receiver)
		server.Logger.Printf("Sent %d bytes (%d lines) of game server logs to %v\n", bytes, len(server.Console.Output), receiver)
	} else {
		server.Logger.Println("Ignoring: Nothing to send.")
	}
}

func (server *Server) DispatchConsoleCommand(receiver *net.UDPAddr, command string) {
	server.Logger.Printf("Received server command '%s' from %v\n", command, receiver)

	if server.IsRunning() {
		server.Logger.Printf("Sending server command '%s' from %v to game server console...\n", command, receiver)
		server.Console.Input <- command

		if command == "quit" {
			server.Command.Wait()
			server.Delete()
		} else {
			// ensure the console replies are printed to as expected
			time.Sleep(serverConsoleInputDelay)

			server.Logger.Printf("Sending game server console replies for command '%s' to %v...\n", command, receiver)
			server.Console.SendCommandReplies(command, receiver)
		}
	}
}
